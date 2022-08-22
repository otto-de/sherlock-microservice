package dockerimage

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

const (
	// remoteDigestPrefix is the line start we search for to get the remote
	// digest
	remoteDigestPrefix = "latest: digest: sha256:"
)

// Push calls docker and pushes a tag.
// Returns the remote digest for the pushed image.
func Push(tag string) (string, error) {

	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()

	digest := make(chan string)
	defer close(digest)
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if !strings.HasPrefix(scanner.Text(), remoteDigestPrefix) {
				continue
			}

			digest <- scanner.Text()[len(remoteDigestPrefix) : len(remoteDigestPrefix)+64]
			return
		}
	}()
	err := executeDockerCommandWithStdout(io.MultiWriter(os.Stdout, w), "push", tag)
	if err != nil {
		return "", err
	}

	select {
	case d := <-digest:
		return d, nil
	default:
		return "", errors.New("No digest in output found")
	}
}
