# Changelog

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
