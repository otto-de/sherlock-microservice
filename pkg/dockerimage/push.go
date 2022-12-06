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
func PushWithWriters(tag string, out io.Writer, errOut io.Writer) (string, error) {
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
	err := executeDockerCommandWithWriters(io.MultiWriter(out, w), errOut, "push", tag)
	if err != nil {
		return "", err
	}

	select {
	case d := <-digest:
		return d, nil
	default:
		return "", errors.New("no digest in output found")
	}
}

// Push calls docker and pushes a tag.
// Returns the remote digest for the pushed image.
func Push(tag string) (string, error) {
	return PushWithWriters(tag, os.Stdout, os.Stderr)
}
