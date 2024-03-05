package dockerimage

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
)

const (
	// remoteDigestPrefix is the line start we search for to get the remote
	// digest
	remoteDigestPrefix = "latest: digest: sha256:"
)

var (
	DigestOutputError = digestOutputError{
		errors.New("no digest in output found"),
	}
)

type digestOutputError struct {
	error
}

func (err digestOutputError) Is(other error) bool {
	_, ok := other.(digestOutputError)
	return ok
}

// Push calls docker and pushes a tag.
// Returns the remote digest for the pushed image.
func PushWithWriters(tag string, out io.Writer, errOut io.Writer) (string, error) {
	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()

	digest := make(chan string, 1)
	errs := make(chan error, 1)
	errG := sync.WaitGroup{}

	errG.Add(1)
	go func() {
		defer close(digest)
		defer errG.Done()

		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if !strings.HasPrefix(scanner.Text(), remoteDigestPrefix) {
				continue
			}

			digest <- scanner.Text()[len(remoteDigestPrefix) : len(remoteDigestPrefix)+64]
			return
		}
		errs <- DigestOutputError
	}()

	errG.Add(1)
	go func() {
		defer errG.Done()

		errs <- executeDockerCommandWithWriters(io.MultiWriter(out, w), errOut, "push", tag)
	}()

	go func() {
		defer close(errs)

		errG.Wait()
	}()

	d, ok := <-digest
	if ok {
		return d, nil
	}
	for err := range errs {
		if err != nil {
			return "", err
		}
	}
	panic("should have resulted in error before")
}

// Push calls docker and pushes a tag.
// Returns the remote digest for the pushed image.
func Push(tag string) (string, error) {
	return PushWithWriters(tag, os.Stdout, os.Stderr)
}
