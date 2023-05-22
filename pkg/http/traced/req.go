package traced

import (
	"context"
	"io"
	"net/http"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go.opentelemetry.io/otel/propagation"
)

// NewIdempotentPostRequestWithContext creates a POST request which is marked as being idempotent.
// Thus it is ok for the Go network transport, Proxies, etc. to redeliver the request payload.
// Also injects CloudTrace information from `Context`.
func NewIdempotentPostRequestWithContext(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	// According to https://github.com/golang/go/commit/bfd9b94069e74b0c6516a045cbb83bf1024a1269
	// this lets the underlying Go transport retry
	r.Header.Set("Idempotency-Key", "x")
	propagator.CloudTraceFormatPropagator{}.Inject(ctx, propagation.HeaderCarrier(r.Header))

	return r, nil
}
