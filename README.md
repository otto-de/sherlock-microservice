[![Go Reference](https://pkg.go.dev/badge/github.com/otto-de/sherlock-microservice.svg)](https://pkg.go.dev/github.com/otto-de/sherlock-microservice)
[![Go Report Card](https://goreportcard.com/badge/github.com/otto-de/sherlock-microservice)](https://goreportcard.com/report/github.com/otto-de/sherlock-microservice)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/otto-de/sherlock-microservice)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/otto-de/sherlock-microservice)

# Microservice module for Sherlock

Contains following Go packages:

| Package                  | Description |
| ---                      | ---         |
| [closer](pkg/closer)     | Handling for io.Closer instances. |
| [datastorehandlers](pkg/datastorehandlers) | The way we interpret Entries in Datastore. |
| [envflags](pkg/envflags) | Extends flags with Environment Variable handling. |
| [gcp](pkg/gcp)           | Combines all Google Cloud service that we use everywhere. |
| [gke](pkg/gke)           | Ease working with Google Kubernetes Engine. |

Contains following Terraform modules:

| Module                                    | Description |
| ---                                       | ---         |
| [service_cluster](tf/service_cluster)     | Module for creating GKE Clusters for Services |
| [service_container](tf/service_container) | Module for building and uploading Container images for our Services |
| [service_namespace](tf/service_namespace) | Module for creating and configuring Kubernetes Namespaces for our Services |
