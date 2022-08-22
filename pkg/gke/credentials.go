package gke

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"k8s.io/client-go/tools/clientcmd/api"
)

// KubeContext the name of Clusters in ~./kube/config
type KubeContext string

// AddAuthForCluster adds credentials for a Cluster to ~./kube/config.
// Currently uses gcloud since it quite involved to do this in pure Go.
func AddAuthForCluster(config api.Config, project, location, clusterName string) (KubeContext, error) {

	cmd := exec.Command("gcloud", "container", "clusters", "get-credentials", clusterName, "--region", location, "--project", project)
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	defer errReader.Close()
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer outReader.Close()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(os.Stderr, errReader)
	}()
	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, outReader)
	}()
	err = cmd.Run()
	wg.Wait()
	kctx := fmt.Sprintf("gke_%s_%s_%s", project, location, clusterName)
	return KubeContext(kctx), err
	/*
		ctx := context.Background()
		containerService, err := container.NewService(ctx)
		if err != nil {
			return "", fmt.Errorf("Creating Container Service failed: %w", err)
		}

		pzcs := container.NewProjectsZonesClustersService(containerService)
		cluster, err := pzcs.Get(project, "", clusterName).Name(fmt.Sprintf("projects/%s/locations/%s/clusters/%s", project, location, clusterName)).Do()
		if err != nil {
			return "", fmt.Errorf("Getting Cluster `%s` failed: %w", clusterName, err)
		}

		//panic(fmt.Sprintf("%#v", config.AuthInfos))
		kctx := fmt.Sprintf("gke_%s_%s_%s", project, location, clusterName)

		auths := config.AuthInfos[kctx]
		if auths == nil {
			auths = api.NewAuthInfo()
			config.AuthInfos[kctx] = auths
		}
		//cluster.MasterAuth.ClusterCaCertificate
		auths.ClientKeyData = []byte(cluster.MasterAuth.ClientKey)
		auths.ClientCertificateData = []byte(cluster.MasterAuth.ClientCertificate)

		err = clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), config, false)
		if err != nil {
			return "", fmt.Errorf("Modifying Config failed: %w", err)
		}

		return KubeContext(kctx), nil
	*/
}

/*
class APIAdapter(object):
  """Handles making api requests in a version-agnostic way."""

  def __init__(self, api_version, registry, client, messages):
    self.api_version = api_version
    self.registry = registry
    self.client = client
    self.messages = messages

  def ParseCluster(self, name):
    properties.VALUES.compute.zone.Get(required=True)
    properties.VALUES.core.project.Get(required=True)
    return self.registry.Parse(
        name, collection='container.projects.zones.clusters')
*/

/*
cluster_ref = adapter.ParseCluster(args.name)

log.status.Print('Fetching cluster endpoint and auth data.')
# Call DescribeCluster to get auth info and cache for next time
cluster = adapter.GetCluster(cluster_ref)
if not adapter.IsRunning(cluster):
  log.error(
	  'cluster %s is not running. The kubernetes API will probably be '
	  'unreachable.' % cluster_ref.clusterId)
util.ClusterConfig.Persist(cluster, cluster_ref.projectId, self.cli)
*/
