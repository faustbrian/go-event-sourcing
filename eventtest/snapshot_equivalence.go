package eventtest

import (
	"context"
	"fmt"

	eventsourcing "github.com/faustbrian/go-event-sourcing/v2"
)

// SnapshotEquivalenceConfig supplies independent full-history and
// snapshot-accelerated aggregate loaders plus application equality.
type SnapshotEquivalenceConfig[Aggregate any] struct {
	// FullHistory is required and independently reconstructs authoritative history
	// using the supplied context.
	FullHistory func(context.Context) (Aggregate, error)
	// Snapshot is required and reports the positive snapshot base actually used; a
	// zero base fails equivalence.
	Snapshot func(context.Context) (Aggregate, uint64, error)
	// Version extracts the observed aggregate version for equality; snapshot
	// checks require it.
	Version func(Aggregate) uint64
	// Equal is the required aggregate-state predicate; versions are compared
	// separately and both aggregate values remain application-owned.
	Equal func(Aggregate, Aggregate) bool
}

// CheckSnapshotEquivalence proves that a load which actually used a snapshot
// reaches the same aggregate state and version as authoritative full replay.
func CheckSnapshotEquivalence[Aggregate any](
	ctx context.Context,
	config SnapshotEquivalenceConfig[Aggregate],
) error {
	if ctx == nil ||
		config.FullHistory == nil ||
		config.Snapshot == nil ||
		config.Version == nil ||
		config.Equal == nil {
		return eventsourcing.ErrInvalidArgument
	}
	full, err := config.FullHistory(ctx)
	if err != nil {
		return fmt.Errorf("load full event history: %w", err)
	}
	accelerated, snapshotVersion, err := config.Snapshot(ctx)
	if err != nil {
		return fmt.Errorf("load snapshot and later history: %w", err)
	}
	fullVersion := config.Version(full)
	acceleratedVersion := config.Version(accelerated)
	if snapshotVersion == 0 {
		return fmt.Errorf("%w: snapshot load did not use a snapshot", ErrConformance)
	}
	if snapshotVersion > acceleratedVersion {
		return fmt.Errorf("%w: snapshot version exceeds aggregate version", ErrConformance)
	}
	if fullVersion != acceleratedVersion {
		return fmt.Errorf("%w: aggregate versions differ", ErrConformance)
	}
	if !config.Equal(full, accelerated) {
		return fmt.Errorf("%w: aggregate states differ", ErrConformance)
	}

	return nil
}
