package test

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/otto-de/sherlock-microservice/pkg/gke"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// PodRun is an test execution run of a Pod.
type PodRun struct {
	ctx          context.Context
	logStreaming sync.WaitGroup
	pods         v1.PodInterface
	Pod          *core.Pod
	streams      genericclioptions.IOStreams
}

// NewPodRun is starting a remote Pod and streams its output locally.
// Output gets written using provided streams.
func NewPodRunWithStreams(clientset *kubernetes.Clientset, ctx context.Context, pod *core.Pod, streams genericclioptions.IOStreams) *PodRun {

	pods := clientset.CoreV1().Pods(pod.Namespace)

	fmt.Fprintln(streams.Out, "Creating Test Pod")
	pod, err := pods.Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(streams.Out, "Created Test Pod `%s`\n", pod.Name)

	pr := &PodRun{
		ctx:     ctx,
		pods:    pods,
		Pod:     pod,
		streams: streams,
	}
	pr.logStreaming.Add(1)
	go func() {
		defer pr.logStreaming.Done()

		err := gke.StreamContainerLog(pod.Namespace, pod, "test", streams)
		if err != nil {
			panic(err)
		}
	}()
	return pr
}

// NewPodRun is starting a remote Pod and streams its output locally.
// Output gets written to Stdout and Stderr.
func NewPodRun(clientset *kubernetes.Clientset, ctx context.Context, pod *core.Pod) *PodRun {

	streams := genericclioptions.IOStreams{
		In:     nil,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	return NewPodRunWithStreams(clientset, ctx, pod, streams)
}

// Close waits until there is no more output to stream. Then deletes the Pod.
func (pr *PodRun) Close() error {
	pr.logStreaming.Wait()
	fmt.Fprintf(pr.streams.Out, "Deleting Test Pod `%s`\n", pr.Pod.Name)
	return pr.pods.Delete(pr.ctx, pr.Pod.Name, metav1.DeleteOptions{})
}
