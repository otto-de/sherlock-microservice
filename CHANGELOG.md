# Changelog

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
