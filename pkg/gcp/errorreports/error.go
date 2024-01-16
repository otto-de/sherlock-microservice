package errorreports

import (
	"cloud.google.com/go/errorreporting"
)

// Error can be used to build a `errorreporting.Entry`.
type Error interface {
	error

	// ReportingEntry uses the saved error information to build
	// an `errorreporting.Entry`.
	ReportingEntry() *errorreporting.Entry
}
