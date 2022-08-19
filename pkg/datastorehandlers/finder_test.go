package datastorehandlers

import (
	"errors"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/google/go-cmp/cmp"
)

func TestFuzzyMatchEmpty(t *testing.T) {
	_, err := fuzzyMatch([]*datastore.Key{},
		234,
		nil,
	)
	if err == nil {
		t.Fatal("Expected error in fuzzyMatch.")
	}
	if !errors.Is(err, datastore.ErrNoSuchEntity) {
		t.Fatal("FuzzyMatch returned unexpected error:", err)
	}
}

func TestFuzzyMatchExact(t *testing.T) {

	m, err := fuzzyMatch([]*datastore.Key{
		{
			Kind: "testkind",
			Name: "1234",
		},
		{
			Kind: "testkind",
			Name: "01235",
		},
		{
			Kind: "testkind",
			Name: "234",
		},
		{
			Kind: "testkind",
			Name: "235",
		},
	},
		234,
		nil,
	)
	if err != nil {
		t.Fatal("Unexpected error in fuzzyMatch:", err)
	}
	d := cmp.Diff(m, &datastore.Key{
		Kind: "testkind",
		Name: "234",
	})
	if d != "" {
		t.Fatal("Mismatch between found key:", d)
	}
}

func TestNearFuzzyMatch(t *testing.T) {
	m, err := fuzzyMatch([]*datastore.Key{
		{
			Kind: "testkind",
			Name: "1234",
		},
		{
			Kind: "testkind",
			Name: "01235",
		},
		{
			Kind: "testkind",
			Name: "235",
		},
	},
		234,
		nil,
	)
	if err != nil {
		t.Fatal("Unexpected error in fuzzyMatch:", err)
	}
	d := cmp.Diff(m, &datastore.Key{
		Kind: "testkind",
		Name: "235",
	})
	if d != "" {
		t.Fatal("Mismatch between found key:", d)
	}
}
