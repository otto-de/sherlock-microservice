package datastorehandlers

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/logging"
)

// FuzzyFinder implements a range find, treating Entry Key Names as epoch values.
// Assumes Name is a epoch in Milliseconds.
type FuzzyFinder struct {
	ancestorKind string
	client       *datastore.Client
	kind         string
	maxDeltaMs   int64
	logger       *logging.Logger
}

// NewFuzzyFinder creates a new FuzzyFinder.
// Enforces a present logging.Client.
func NewFuzzyFinder(client *datastore.Client, ancestorKind, kind string, maxDeltaMs int64, l *logging.Client) *FuzzyFinder {
	var logger *logging.Logger
	if l != nil {
		logger = l.Logger("FuzzyFinder")
	}
	return &FuzzyFinder{
		ancestorKind: ancestorKind,
		client:       client,
		kind:         kind,
		maxDeltaMs:   maxDeltaMs,
		logger:       logger,
	}
}

// Close flushes potential logging messages
func (f *FuzzyFinder) Close() {
	if f.logger != nil {
		f.logger.Flush()
	}
}

func (f *FuzzyFinder) FindClosestToEpochMs(ctx context.Context, ancestorName string, epochMs int64) (*datastore.Key, error) {

	lowerName := fmt.Sprintf("%013d", epochMs-f.maxDeltaMs)
	upperName := fmt.Sprintf("%013d", epochMs+f.maxDeltaMs)

	// Due to limitations within Datastore we are forced to convert int64 to string.
	// To make all this easire to handle, we only support epochs around time.Now()
	// Be explicit about not supporting too divergent epoch values.
	if len(lowerName) != 13 {
		return nil, fmt.Errorf("not supported value due to digit count (%d): %s", len(lowerName), lowerName)
	}

	if len(upperName) != 13 {
		return nil, fmt.Errorf("not supported value due to digit count (%d): %s", len(upperName), upperName)
	}

	ancestor := datastore.NameKey(f.ancestorKind, ancestorName, nil)
	lower := datastore.NameKey(f.kind, lowerName, ancestor)
	upper := datastore.NameKey(f.kind, upperName, ancestor)
	q := datastore.NewQuery(f.kind).Ancestor(ancestor).Filter("__key__ >", lower).Filter("__key__ <", upper).Order("__key__").KeysOnly()

	foundKeys, err := f.client.GetAll(ctx, q, nil)
	if err != nil {
		return nil, err
	}

	key, err := fuzzyMatch(foundKeys, epochMs, f.logger)
	if key == nil && len(foundKeys) != 0 && f.logger != nil {
		f.logger.Log(logging.Entry{
			Severity: logging.Debug,
			Payload:  fmt.Sprintf("Fuzzy matching matched no key of: %s", foundKeys),
		})
	}
	return key, err
}

// fuzzyMatch searches `datastore.Key` `Name` that is numerically closest to `epochMs`.
// No search stability is guaranteed!
// Currently does not expect keys to be ordered.
// Has O(n) time complexity.
func fuzzyMatch(keys []*datastore.Key, epochMs int64, logger *logging.Logger) (*datastore.Key, error) {
	nearestI := -1
	var nearestDeltaMs int64 = math.MaxInt64
	for i, foundKey := range keys {
		foundEpochMs, err := strconv.ParseInt(foundKey.Name, 10, 64)
		if err != nil {
			// Seems like datastore is corrupt
			logger.Log(logging.Entry{
				Severity: logging.Alert,
				Payload:  fmt.Sprintf("Invalid key name in database: %s", foundKey.Name),
			})
			continue
		}

		deltaMs := foundEpochMs - epochMs
		if deltaMs < 0 {
			deltaMs = -deltaMs
		}

		if deltaMs < nearestDeltaMs {
			nearestDeltaMs = deltaMs
			nearestI = i
		}
	}

	if nearestI == -1 {
		return nil, datastore.ErrNoSuchEntity
	}

	return keys[nearestI], nil
}
