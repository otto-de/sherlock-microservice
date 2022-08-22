package dockerimage

import (
	"io"
	"os"
	"os/exec"
	"sync"
)

// executeDockerCommandWithStdout calls docker CLI with `arg`.
// Ensures that output is written to provided `io.Writer`.
func executeDockerCommandWithStdout(stdout io.Writer, arg ...string) error {
	cmd := exec.Command("docker", arg...)
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer errReader.Close()
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
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
		io.Copy(stdout, outReader)
	}()
	err = cmd.Run()
	wg.Wait()
	return err
}

// executeDockerCommand call docker CLI with `arg`.
func executeDockerCommand(arg ...string) error {
	return executeDockerCommandWithStdout(os.Stdout, arg...)
}
