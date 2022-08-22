package dockerimage

import "os"

// BuildFromContainerfile builds a Container image from a Containerfile.
// Invokes docker via CLI since Docker's dependencies and API are beyond brittle.
func BuildFromContainerfile(path string, tag string) error {

	dir, _ := os.Getwd()
	os.Chdir(path)
	defer os.Chdir(dir)
	// Normally you would want to use a Golang API here.
	// Since the Docker API is a very thin shim, badly documented
	// and talking to Dockerd or Buildkit has a lot of quirks
	// just call the binary here and be done with it.
	return executeDockerCommand("build", "-f", "Containerfile", ".", "-t", tag)
}
