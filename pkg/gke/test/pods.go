package gke

import (
	"context"
	"sync"
	"testing"

	"github.com/otto-de/sherlock-microservice/pkg/gke"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// PodRun is an test execution run of a Pod.
type PodRun struct {
	ctx          context.Context
	logStreaming sync.WaitGroup
	pods         v1.PodInterface
	Pod          *core.Pod
	t            *testing.T
}

// NewPodRun is starting a remote Pod and streams its output locally.
func NewPodRun(t *testing.T, clientset *kubernetes.Clientset, ctx context.Context, pod *core.Pod) *PodRun {
	pods := clientset.CoreV1().Pods(pod.Namespace)

	t.Log("Creating Test Pod")
	pod, err := pods.Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	t.Logf("Created Test Pod `%s`", pod.Name)

	pr := &PodRun{
		ctx:  ctx,
		pods: pods,
		Pod:  pod,
		t:    t,
	}
	pr.logStreaming.Add(1)
	go func() {
		defer pr.logStreaming.Done()
		err := gke.StreamContainerLog(pod.Namespace, pod, "test")
		if err != nil {
			panic(err)
		}
	}()
	return pr
}

// Close waits until there is no more output to stream. Then deletes the Pod.
func (pr *PodRun) Close() error {
	pr.logStreaming.Wait()
	pr.t.Logf("Deleting Test Pod `%s`", pr.Pod.Name)
	return pr.pods.Delete(pr.ctx, pr.Pod.Name, metav1.DeleteOptions{})
}
