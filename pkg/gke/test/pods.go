package test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

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
func NewPodRunWithStreams(tb testing.TB, clientset *kubernetes.Clientset, ctx context.Context, pod *core.Pod, streams genericclioptions.IOStreams) *PodRun {
	pods := clientset.CoreV1().Pods(pod.Namespace)

	tb.Log("Creating Test Pod\n")
	pod, err := pods.Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	tb.Logf("Created Test Pod `%s`\n", pod.Name)

	pr := &PodRun{
		ctx:     ctx,
		pods:    pods,
		Pod:     pod,
		streams: streams,
	}
	pr.logStreaming.Add(1)
	go func() {
		defer pr.logStreaming.Done()

		err := gke.StreamContainerLog(ctx, pods, pod, "test", streams)
		if err != nil {
			panic(err)
		}
	}()
	return pr
}

// NewPodRun is starting a remote Pod and streams its output locally.
// Output gets written to Stdout and Stderr.
func NewPodRun(tb testing.TB, clientset *kubernetes.Clientset, ctx context.Context, pod *core.Pod) *PodRun {

	streams := genericclioptions.IOStreams{
		In:     nil,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	return NewPodRunWithStreams(tb, clientset, ctx, pod, streams)
}

// Close waits until there is no more output to stream. Then deletes the Pod.
func (pr *PodRun) Close(tb testing.TB, clientset *kubernetes.Clientset, ctx context.Context) error {
	pr.logStreaming.Wait()

	pods := clientset.CoreV1().Pods(pr.Pod.Namespace)

	pod, err := pods.Get(ctx, pr.Pod.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get Pod '%s': %w", pr.Pod.Name, err)
	}

	switch pod.Status.Phase {
	case core.PodFailed:
		tb.Logf("Pod '%s' failed. Keeping failed Test Pod for debugging\n", pod.Name)
		return nil
	case core.PodSucceeded:
		tb.Logf("Pod '%s' succeeded. Deleting Test Pod\n", pod.Name)
		return pods.Delete(ctx, pod.Name, metav1.DeleteOptions{})
	default:
		tb.Logf("Pod '%s' in state '%v'. Keeping Test Pod for possible debugging\n", pod.Name, pod.Status.Phase)
		return nil
	}
}
