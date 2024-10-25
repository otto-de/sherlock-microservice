package gke

import (
	"context"
	"io"
	"time"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// StreamContainerLog streams log for that one container to Stdout and Stderr.
func StreamContainerLog(ctx context.Context, pods v1.PodInterface, pod *core.Pod, containerName string, streams genericclioptions.IOStreams) error {
	// wait for pod to be running or terminated
	// otherwise we cant get logs
	for {
		spod, err := pods.Get(ctx, pod.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if spod.Status.Phase == core.PodRunning || spod.Status.Phase == core.PodSucceeded || spod.Status.Phase == core.PodFailed {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	r := pods.GetLogs(pod.Name, &core.PodLogOptions{Follow: true, Container: containerName})
	rc, err := r.Stream(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()
	_, err = io.Copy(streams.Out, rc)
	if err != nil {
		return err
	}
	return nil
}
