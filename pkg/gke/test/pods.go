package test

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/otto-de/sherlock-microservice/pkg/gke"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type CreatedPod struct {
	pods         v1.PodInterface
	pod          *core.Pod
	logStreaming sync.WaitGroup
}

func MustCreatePod(tb testing.TB, clientset *kubernetes.Clientset, ctx context.Context, pod *core.Pod) *CreatedPod {
	pods := clientset.CoreV1().Pods(pod.Namespace)

	tb.Log("Creating Test Pod\n")
	pod, err := pods.Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	tb.Logf("Created Test Pod `%s`\n", pod.Name)

	return &CreatedPod{
		pods: pods,
		pod:  pod,
	}
}

func (cp *CreatedPod) Delete() error {
	cp.logStreaming.Wait()

	return cp.pods.Delete(context.TODO(), cp.pod.Name, metav1.DeleteOptions{})
}

// PodRun is an test execution run of a Pod.
type PodRun struct {
	ctx     context.Context
	Pod     *CreatedPod
	streams genericiooptions.IOStreams
}

// RunWithStreams streams a Pod output locally.
// Output gets written using provided streams.
func (cp *CreatedPod) RunWithStreams(tb testing.TB, ctx context.Context, streams genericiooptions.IOStreams) *PodRun {

	pr := &PodRun{
		ctx:     ctx,
		Pod:     cp,
		streams: streams,
	}
	cp.logStreaming.Add(1)
	go func() {
		defer cp.logStreaming.Done()

		err := gke.StreamContainerLog(ctx, cp.pods, cp.pod, "test", streams)
		if err != nil {
			panic(err)
		}
	}()
	return pr
}

// Run streams a Pod output locally.
// Output gets written to Stdout and Stderr.
func (cp *CreatedPod) Run(tb testing.TB, clientset *kubernetes.Clientset, ctx context.Context) *PodRun {

	streams := genericiooptions.IOStreams{
		In:     nil,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	return cp.RunWithStreams(tb, ctx, streams)
}

// Close waits until there is no more output to stream.
func (pr *PodRun) Close() error {
	pr.Pod.logStreaming.Wait()

	return nil
}
