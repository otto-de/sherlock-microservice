package gke

import (
	"fmt"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

func MonitoredResource(l *logging.Client, project, clusterName, namespace, pod, containerName string) *monitoredres.MonitoredResource {
	instanceId := "unknown"
	zone := "unknown"
	res := &monitoredres.MonitoredResource{
		Type: "gke_container",
		Labels: map[string]string{
			"project_id":     project,
			"cluster_name":   clusterName,
			"namespace_id":   namespace,
			"instance_id":    instanceId,
			"pod_id":         pod,
			"container_name": containerName,
			"zone":           zone,
		},
	}

	logger := l.Logger("MonitoredResource")
	defer logger.Flush()

	// Refine labels
	var err error
	instanceId, err = metadata.InstanceID()
	if instanceId == "" {
		logger.Log(logging.Entry{
			Severity: logging.Info,
			Payload:  fmt.Sprintf("Error getting instance ID: %s", err),
			Resource: res,
		})
	} else {
		res.Labels["instance_id"] = instanceId
	}

	zone, err = metadata.Zone()
	if zone == "" {
		logger.Log(logging.Entry{
			Severity: logging.Info,
			Payload:  fmt.Sprintf("Error getting zone: %s", err),
			Resource: res,
		})
	} else {
		res.Labels["zone"] = zone
	}

	return res
}
