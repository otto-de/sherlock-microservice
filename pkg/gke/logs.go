package gke

import (
	"strings"
	"time"

	core "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/logs"
	"k8s.io/kubectl/pkg/polymorphichelpers"
)

var defaultConfigFlags = genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)

// StreamContainerLog streams log for that one container to Stdout and Stderr.
func StreamContainerLog(namespace string, pod *core.Pod, containerName string, streams genericclioptions.IOStreams) error {

	lo := logs.NewLogsOptions(streams, false)
	lo.Follow = true
	lo.Container = containerName
	lo.IgnoreLogErrors = true
	lo.Namespace = namespace
	lo.ConsumeRequestFn = logs.DefaultConsumeRequest
	lo.RESTClientGetter = defaultConfigFlags
	lo.Object = pod
	lo.LogsForObject = polymorphichelpers.LogsForObjectFn
	var err error
	lo.Options, err = lo.ToLogOptions()
	if err != nil {
		return err
	}
	err = lo.Validate()
	for err == nil {
		err = lo.RunLogs()
		if err == nil {
			return err
		}
		if strings.HasSuffix(err.Error(), "is waiting to start: ContainerCreating") {
			err = nil
			time.Sleep(time.Millisecond * 100)
		}
	}

	return err
}
