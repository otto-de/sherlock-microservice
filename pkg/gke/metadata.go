package gke

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"cloud.google.com/go/compute/metadata"
)

type GKEMetaData struct {
	ClusterLocation  string
	ClusterName      string
	ContainerName    string
	Namespace        string
	InstanceName     string
	InstanceID       int64
	NumericProjectID int64
	PodName          string
	ProjectID        string
	Zone             string
}

type instance struct {
	Instance struct {
		Attributes struct {
			ClusterLocation string `json:"cluster-location"`
			ClusterName     string `json:"cluster-name"`
			ClusterUID      string `json:"cluster-uid"`
		} `json:"attributes"`
		Hostname          string `json:"hostname"`
		ID                int64  `json:"id"`
		Name              string `json:"name"`
		NetworkInterfaces map[string]struct {
			Ipv6s string `json:"ipv6s"`
		} `json:"networkInterfaces"`
		ServiceAccounts map[string]struct {
			Aliases []string `json:"aliases"`
			Email   string   `json:"email"`
			Scopes  []string `json:"scopes"`
		} `json:"serviceAccounts"`
		Zone string `json:"zone"`
	} `json:"instance"`
	Project struct {
		NumericProjectID int64  `json:"numericProjectId"`
		ProjectID        string `json:"projectId"`
	} `json:"project"`
}

func parseInstanceJSON(jsonData []byte) (*instance, error) {
	var instance instance
	err := json.Unmarshal(jsonData, &instance)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func GetMetaData(ctx context.Context) (*GKEMetaData, error) {
	s, err := metadata.GetWithContext(ctx, "/?recursive=true")
	if err != nil {
		return &GKEMetaData{}, err
	}

	i, err := parseInstanceJSON([]byte(s))
	if err != nil {
		return &GKEMetaData{}, err
	}

	zp := strings.Split(i.Instance.Zone, "/")
	zoneName := zp[len(zp)-1]

	return &GKEMetaData{
		ClusterLocation:  i.Instance.Attributes.ClusterLocation,
		ClusterName:      i.Instance.Attributes.ClusterName,
		ContainerName:    os.Getenv("CONTAINER_NAME"),
		Namespace:        os.Getenv("POD_NAMESPACE"),
		InstanceName:     i.Instance.Name,
		InstanceID:       i.Instance.ID,
		NumericProjectID: i.Project.NumericProjectID,
		PodName:          os.Getenv("POD_NAME"),
		ProjectID:        i.Project.ProjectID,
		Zone:             zoneName,
	}, nil
}
