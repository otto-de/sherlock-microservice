package closer

import (
	"fmt"
	"io"

	"cloud.google.com/go/logging"
)

// CloseOrLog should be used in cases where it is ok
// to not handle failures in Close calls.
// The function then logs a warning to Cloud Logging.
// TODO: Add option to make log entry more concrete.
func CloseOrLog(l *logging.Logger, c io.Closer) {
	err := c.Close()
	if err != nil {
		l.Log(logging.Entry{
			Severity: logging.Warning,
			Payload:  fmt.Sprintf("Closing %#v failed: %s", c, err),
		})
	}
}
