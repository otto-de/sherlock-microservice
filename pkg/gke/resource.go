package gke

import (
	"context"
	"fmt"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

func MonitoredResource(ctx context.Context, l *logging.Client, project, clusterName, namespace, pod, containerName string) *monitoredres.MonitoredResource {
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
	instanceId, err = metadata.InstanceIDWithContext(ctx)
	if instanceId == "" {
		logger.Log(logging.Entry{
			Severity: logging.Info,
			Payload:  fmt.Sprintf("Error getting instance ID: %s", err),
			Resource: res,
		})
	} else {
		res.Labels["instance_id"] = instanceId
	}

	zone, err = metadata.ZoneWithContext(ctx)
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
			"instance_id":    fmt.Sprint(metadata.InstanceID),
			"pod_id":         metadata.PodName,
			"container_name": metadata.ContainerName,
			"zone":           metadata.Zone,
		},
	}
}

func TraceResourceFromMetaData(
	serviceName string,
	metadata *GKEMetaData,
) *resource.Resource {
	return resource.NewSchemaless(
		semconv.CloudProviderGCP,
		semconv.CloudPlatformGCPKubernetesEngine,
		semconv.CloudAccountID(metadata.ProjectID),
		semconv.CloudAvailabilityZone(metadata.Zone),
		semconv.ServiceName(serviceName),
		semconv.K8SClusterName(metadata.ClusterName),
		semconv.K8SNamespaceName(metadata.Namespace),
		semconv.K8SNodeName(metadata.InstanceName),
		semconv.K8SPodName(metadata.PodName),
		semconv.K8SContainerName(metadata.ContainerName),
		attribute.String("gcp.project_id", metadata.ProjectID),
	)
}
