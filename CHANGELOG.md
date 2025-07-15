# Changelog

## v0.0.80 Change DiscoverServices behavior
It was discovered that the runtime of service discovery is very lengthy in some scenarios.
As such now change the API so that the API users are forced to handle this - in an hopefully better way.
- DiscoverServices changed so that it starts discovery immediately but returns the result via channel
- DiscoverServicesOnce allows calling discovery concurrently

## v0.0.66 Allow for caching in aws_sns_message_verify

## v0.0.65 Change pod log streaming
- respect the connection settings of the clientset

## v0.0.64 Implement Benthos processor aws_sns_message_verify
- Allow for verifying SNS Webhook data

## v0.0.63 Enforce Docker file paths to building Containers

## v0.0.62 Switch GKE Pods to execute on TB
- Interacting with GKE Pod functions new needs testing.TB object.

## v0.0.60 Allow TraceExporterOption
- Enable passing options to `DiscoverServices` via `WithTraceExporterOption`

## v0.0.55 Allow setting of Extensions
- Introduces `ApplyCloudEventOptions` for applying Options to CloudEvents

## v0.0.51 Allow passing of Test Clients
- Introduces `FakeSetup`, which allows for passing TestClients along.

## v0.0.47 Fix accessing OrderingKey
**Bugfixes**
- Corrects implementation of function `ApplyCloudEventsPubSubOrderingKey`

## v0.0.45 Accessing OrderingKey from callsites
- Add function `ApplyCloudEventsPubSubOrderingKey`

## v0.0.44 ErrorReporting
- Add interface `errorreports.Error`

  Special error interface so that it is possible to generate a GCP ErrorReport directly from `error`.
- Add interface `publisher.Publisher`

  Use `Publisher` for rich publish integration. Resulting error can directly be used for GCP ErrorReporting.

## v0.0.39 More options for envflags
- Add function `envflags.GetBoolDefault`

## v0.0.38 Further options for ErrorReporting
- Add WithErrorReportCallback for testing Services.

## v0.0.36 Ease for ErrorReporting
- Add WithErrorReportChannel for testing Services.

## v0.0.34 Ease testing of PubSub
- Add function `gcp.pubsubcmp.DiffMessages`.

## v0.0.5 - Migrate `dockerimage` to _sherlock-mavrodi_

- Add function `dockerimage.BuildFromContainerfile` for building Container Image.
- Add function `dockerimage.Push` for pushing a Container Image.

**Bugfixes**
- Correct docs of `StreamContainerLog`


---
## v0.0.4 - Migrate `gke` package to _sherlock-mavrodi_

- Add function `gke.AddAuthForCluster` for having local config.
- Add function `gke.StreamContainerLog` for delegating a containers remote log to local files.
- Add new type `gke.test.PodRun` for executing a Pod for testing.
