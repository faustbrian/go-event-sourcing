package eventsourcing

import (
	"bytes"
	"context"
	"maps"
	"time"
)

const (
	// MaxSnapshotStateBytes bounds one encoded aggregate snapshot.
	MaxSnapshotStateBytes = 8 << 20
)

// SnapshotInput supplies one immutable derived aggregate snapshot.
type SnapshotInput struct {
	// Stream identifies the aggregate whose derived state is represented.
	Stream StreamID
	// AggregateVersion is the positive last authoritative stream version
	// represented by State.
	AggregateVersion uint64
	// SchemaVersion is the positive snapshot-state schema identity, independent of
	// event schema.
	SchemaVersion SchemaVersion
	// State is nonempty and at most MaxSnapshotStateBytes; NewSnapshot copies it.
	State []byte
	// Metadata is validated/copy-owned application snapshot metadata, subject to
	// core metadata bounds.
	Metadata map[string]string
	// CreatedAt is required and normalized to UTC microseconds.
	CreatedAt time.Time
}

// Snapshot is immutable-by-contract derived aggregate state.
//
// Snapshots accelerate restoration and are never authoritative history.
type Snapshot struct {
	stream           StreamID
	aggregateVersion uint64
	schemaVersion    SchemaVersion
	state            []byte
	metadata         map[string]string
	createdAt        time.Time
}

// NewSnapshot validates and owns derived aggregate state.
func NewSnapshot(input SnapshotInput) (Snapshot, error) {
	if input.Stream.IsZero() {
		return Snapshot{}, invalid("stream", "must be assigned")
	}
	if input.AggregateVersion == 0 {
		return Snapshot{}, invalid("aggregate_version", "must be greater than zero")
	}
	if input.SchemaVersion == 0 {
		return Snapshot{}, invalid("snapshot_schema_version", "must be greater than zero")
	}
	if len(input.State) == 0 || len(input.State) > MaxSnapshotStateBytes {
		return Snapshot{}, invalid("snapshot_state", "must be non-empty and bounded")
	}
	metadata, err := copyMetadata(input.Metadata)
	if err != nil {
		return Snapshot{}, err
	}
	if input.CreatedAt.IsZero() {
		return Snapshot{}, invalid("created_at", "must be assigned")
	}

	return Snapshot{
		stream:           input.Stream,
		aggregateVersion: input.AggregateVersion,
		schemaVersion:    input.SchemaVersion,
		state:            cloneBytes(input.State),
		metadata:         metadata,
		createdAt:        normalizeTime(input.CreatedAt),
	}, nil
}

// Stream returns the aggregate root identity represented by the snapshot.
func (snapshot Snapshot) Stream() StreamID {
	return snapshot.stream
}

// AggregateVersion returns the last stored stream version represented by
// state.
func (snapshot Snapshot) AggregateVersion() uint64 {
	return snapshot.aggregateVersion
}

// SchemaVersion returns the explicit snapshot-state schema version.
func (snapshot Snapshot) SchemaVersion() SchemaVersion {
	return snapshot.schemaVersion
}

// State returns a defensive copy of encoded aggregate state.
func (snapshot Snapshot) State() []byte {
	return cloneBytes(snapshot.state)
}

// Metadata returns a defensive copy of application snapshot metadata.
func (snapshot Snapshot) Metadata() map[string]string {
	return cloneMetadata(snapshot.metadata)
}

// CreatedAt returns the canonical snapshot creation time.
func (snapshot Snapshot) CreatedAt() time.Time {
	return snapshot.createdAt
}

// IsZero reports whether the snapshot has not been assigned.
func (snapshot Snapshot) IsZero() bool {
	return snapshot.stream.IsZero() &&
		snapshot.aggregateVersion == 0 &&
		snapshot.schemaVersion == 0 &&
		len(snapshot.state) == 0 &&
		len(snapshot.metadata) == 0 &&
		snapshot.createdAt.IsZero()
}

// Equal compares every observable snapshot field.
func (snapshot Snapshot) Equal(other Snapshot) bool {
	return snapshot.stream == other.stream &&
		snapshot.aggregateVersion == other.aggregateVersion &&
		snapshot.schemaVersion == other.schemaVersion &&
		bytes.Equal(snapshot.state, other.state) &&
		maps.Equal(snapshot.metadata, other.metadata) &&
		snapshot.createdAt.Equal(other.createdAt)
}

// SnapshotStore persists optional replaceable derived aggregate state.
//
// Save must be atomic. It must reject aggregate or schema regressions, accept
// an exact idempotent retry, and reject different state at the same versions.
// Delete is idempotent.
type SnapshotStore interface {
	// Load returns independently owned derived state or ErrSnapshotNotFound;
	// history remains authoritative.
	Load(context.Context, StreamID) (Snapshot, error)
	// Save atomically rejects version regressions and conflicting equal-version
	// state, while accepting an exact idempotent retry.
	Save(context.Context, Snapshot) error
	// Delete removes derived state idempotently, without modifying event history.
	Delete(context.Context, StreamID) error
}

// SnapshotVersionError describes a rejected aggregate or schema regression.
//
// Error redacts Stream because aggregate identifiers may be sensitive.
type SnapshotVersionError struct {
	// Stream identifies the aggregate whose derived state is represented.
	Stream StreamID
	// StoredAggregateVersion is the retained stream version used for the
	// regression check.
	StoredAggregateVersion uint64
	// IncomingAggregateVersion is the proposed stream version that could not
	// replace retained state.
	IncomingAggregateVersion uint64
	// StoredSchemaVersion is the retained snapshot schema version.
	StoredSchemaVersion SchemaVersion
	// IncomingSchemaVersion is the proposed schema version accompanying the
	// rejected save; either aggregate or schema version may have regressed.
	IncomingSchemaVersion SchemaVersion
}

// Error implements error without disclosing stream identity.
func (*SnapshotVersionError) Error() string {
	return ErrSnapshotStale.Error()
}

// Unwrap classifies the error as ErrSnapshotStale.
func (*SnapshotVersionError) Unwrap() error {
	return ErrSnapshotStale
}

// SnapshotConflictError describes different derived state at identical
// aggregate and schema versions.
//
// Error redacts Stream because aggregate identifiers may be sensitive.
type SnapshotConflictError struct {
	// Stream identifies the aggregate whose derived state is represented.
	Stream StreamID
	// AggregateVersion is the shared positive stream version of the conflicting
	// snapshots.
	AggregateVersion uint64
	// SchemaVersion is the positive snapshot-state schema identity, independent of
	// event schema.
	SchemaVersion SchemaVersion
}

// Error implements error without disclosing stream identity.
func (*SnapshotConflictError) Error() string {
	return ErrSnapshotConflict.Error()
}

// Unwrap classifies the error as ErrSnapshotConflict.
func (*SnapshotConflictError) Unwrap() error {
	return ErrSnapshotConflict
}

// SnapshotRestorationError redacts an application snapshot decoder diagnostic
// while preserving its cause for inspection.
type SnapshotRestorationError struct {
	// Cause retains the underlying failure for Unwrap; Error omits its diagnostic
	// text.
	Cause error
}

// Error implements error without disclosing snapshot state or codec details.
func (*SnapshotRestorationError) Error() string {
	return "aggregate snapshot restoration failed"
}

// Unwrap preserves the decoder cause for errors.Is and errors.As.
func (err *SnapshotRestorationError) Unwrap() error {
	return err.Cause
}

var (
	_ error = (*SnapshotVersionError)(nil)
	_ error = (*SnapshotConflictError)(nil)
	_ error = (*SnapshotRestorationError)(nil)
)
