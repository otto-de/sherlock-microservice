package gke

import (
	"fmt"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
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

func MonitoredResourceFromMetaData(metadata *GKEMetaData) *monitoredres.MonitoredResource {
	return &monitoredres.MonitoredResource{
		Type: "gke_container",
		Labels: map[string]string{
			"project_id":     metadata.ProjectID,
			"cluster_name":   metadata.ClusterName,
			"namespace_id":   metadata.Namespace,
			"instance_id":    string(metadata.InstanceID),
			"pod_id":         metadata.PodName,
			"container_name": metadata.ContainerName,
			"zone":           metadata.Zone,
		},
	}
}

func TraceResourceFromMetaData(serviceName string, metadata *GKEMetaData) *resource.Resource {
	return resource.NewWithAttributes(
		"",
		semconv.CloudProviderGCP,
		semconv.CloudPlatformGCPKubernetesEngine,
		semconv.ServiceNameKey.String(serviceName),
		attribute.String("g.co/r/k8s_container/project_id", metadata.ProjectID),
		attribute.String("g.co/r/k8s_container/cluster_name", metadata.ClusterName),
		attribute.String("g.co/r/k8s_container/namespace", metadata.Namespace),
		attribute.String("g.co/r/k8s_container/location", metadata.ClusterLocation),
		attribute.String("g.co/r/k8s_container/node_name", metadata.InstanceName),
		attribute.String("g.co/r/k8s_container/pod_name", metadata.PodName),
		attribute.String("g.co/r/k8s_container/container_name", metadata.ContainerName),
		attribute.String("g.co/r/k8s_container/zone", metadata.Zone),
	)
}
