package pubsubcmp

import (
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// DiffMessages compares two Messages for equality using all exported Fields.
// Ignores differences in:
// - Attributes.ce-id
// - ID
// - PublishTime
func DiffMessages(lhs pubsub.Message, rhs pubsub.Message) string {
	lhs.ID = ""
	rhs.ID = ""
	lhsAttrs := make(map[string]string, len(lhs.Attributes))
	for k, v := range lhs.Attributes {
		if k == "ce-id" {
			continue
		}
		lhsAttrs[k] = v
	}
	lhs.Attributes = lhsAttrs

	rhsAttrs := make(map[string]string, len(rhs.Attributes))
	for k, v := range rhs.Attributes {
		if k == "ce-id" {
			continue
		}
		rhsAttrs[k] = v
	}
	rhs.Attributes = rhsAttrs

	pt := time.Now()
	lhs.PublishTime = pt
	rhs.PublishTime = pt

	return cmp.Diff(lhs, rhs, cmpopts.IgnoreUnexported(pubsub.Message{}))
}
