package dockerimage

import (
	"io"
	"os/exec"
	"sync"
)

// executeDockerCommandWithWriters calls docker CLI with `arg`.
// Ensures that output is written to provided `io.Writer`s.
func executeDockerCommandWithWriters(out io.Writer, errOut io.Writer, arg ...string) error {
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
		io.Copy(errOut, errReader)
	}()
	go func() {
		defer wg.Done()
		io.Copy(out, outReader)
	}()
	err = cmd.Run()
	wg.Wait()
	return err
}
