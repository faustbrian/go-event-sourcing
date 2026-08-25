package eventsourcing

import (
	"errors"
	"testing"
)

func TestLifecycleGenerationAdvancesExactlyOncePerSuccessfulTransition(t *testing.T) {
	t.Parallel()

	event := invariantDecodedEvent()
	apply := func(DecodedEvent) error { return nil }

	var recorded Lifecycle
	if err := recorded.Record(event, apply); err != nil || recorded.generation != 1 {
		t.Fatalf("Record() error = %v, generation = %d", err, recorded.generation)
	}

	var reconstituted Lifecycle
	if err := reconstituted.Reconstitute(0, nil, apply); err != nil ||
		reconstituted.generation != 1 {
		t.Fatalf("Reconstitute() error = %v, generation = %d", err, reconstituted.generation)
	}

	var restored Lifecycle
	if err := restored.RestoreSnapshotVersion(4); err != nil || restored.generation != 1 {
		t.Fatalf("RestoreSnapshotVersion() error = %v, generation = %d", err, restored.generation)
	}

	var incremental Lifecycle
	if err := incremental.ReconstituteNext(1, nil, apply); err != nil ||
		incremental.generation != 1 {
		t.Fatalf("ReconstituteNext() error = %v, generation = %d", err, incremental.generation)
	}

	acknowledged, changes, prepared, messages := acknowledgementFixture()
	if err := acknowledged.Acknowledge(changes, prepared, messages); err != nil ||
		acknowledged.generation != 8 {
		t.Fatalf("Acknowledge() error = %v, generation = %d", err, acknowledged.generation)
	}
}

func TestLifecycleOverflowAndSegmentBoundsAreExact(t *testing.T) {
	t.Parallel()

	event := invariantDecodedEvent()
	full := Lifecycle{
		committed: ^uint64(0) - 1,
		pending:   []DecodedEvent{event},
	}
	if err := full.Record(event, func(DecodedEvent) error { return nil }); !errors.Is(
		err,
		ErrVersionOverflow,
	) {
		t.Fatalf("Record(at capacity) error = %v", err)
	}

	if _, err := NewHistoricalEvent(HistoricalEventInput{
		SourceVersion: 1,
		SegmentIndex:  MaxUpcastSegments - 1,
		SegmentCount:  MaxUpcastSegments,
		Event:         event,
	}); err != nil {
		t.Fatalf("NewHistoricalEvent(maximum segment) error = %v", err)
	}
	if _, err := NewHistoricalEvent(HistoricalEventInput{
		SourceVersion: 1,
		SegmentCount:  MaxUpcastSegments + 1,
		Event:         event,
	}); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("NewHistoricalEvent(excessive segments) error = %v", err)
	}
}

func TestRestoreSnapshotVersionRejectsEveryPriorStateIndependently(t *testing.T) {
	t.Parallel()

	event := invariantDecodedEvent()
	cases := []Lifecycle{
		{transition: true},
		{committed: 1},
		{pending: []DecodedEvent{event}},
		{generation: 1},
	}
	for index := range cases {
		if err := cases[index].RestoreSnapshotVersion(2); !errors.Is(
			err,
			ErrInvalidLifecycleState,
		) {
			t.Fatalf("RestoreSnapshotVersion(state %d) error = %v", index, err)
		}
	}
}

func TestChangeSetOwnershipChecksEveryIdentityComponent(t *testing.T) {
	t.Parallel()

	for index := 0; index < 5; index++ {
		lifecycle, changes := lifecycleChangeSetFixture()
		switch index {
		case 0:
			changes.owner = &Lifecycle{}
		case 1:
			changes.generation--
		case 2:
			changes.base--
		case 3:
			changes.events = nil
		case 4:
			changes.events = append(changes.events, invariantDecodedEvent())
		}
		if err := lifecycle.MarkPersistenceUnknown(changes); !errors.Is(
			err,
			ErrInvalidChangeSet,
		) {
			t.Fatalf("MarkPersistenceUnknown(mismatch %d) error = %v", index, err)
		}
		if lifecycle.poisoned {
			t.Fatalf("change-set mismatch %d poisoned lifecycle", index)
		}
	}
}

func TestAcknowledgeChecksEveryChangeSetIdentityComponent(t *testing.T) {
	t.Parallel()

	for index := 0; index < 5; index++ {
		lifecycle, changes, prepared, messages := acknowledgementFixture()
		switch index {
		case 0:
			changes.owner = &Lifecycle{}
		case 1:
			changes.generation--
		case 2:
			changes.base--
		case 3:
			changes.events = nil
		case 4:
			changes.events = append(changes.events, invariantDecodedEvent())
		}
		if err := lifecycle.Acknowledge(changes, prepared, messages); !errors.Is(
			err,
			ErrInvalidChangeSet,
		) {
			t.Fatalf("Acknowledge(change-set mismatch %d) error = %v", index, err)
		}
	}
}

func TestAcknowledgeDetectsEachPersistenceMismatchIndependently(t *testing.T) {
	t.Parallel()

	for index := 0; index < 7; index++ {
		lifecycle, changes, prepared, messages := acknowledgementFixture()
		switch index {
		case 0:
			prepared = nil
		case 1:
			messages = nil
		case 2:
			prepared[0].event.name = EventName{value: "account.closed"}
			messages[0].pending = prepared[0]
		case 3:
			prepared[0].event.version++
			messages[0].pending = prepared[0]
		case 4:
			messages[0].pending.id = MessageID{value: "different-message"}
		case 5:
			messages[0].streamVersion++
		case 6:
			messages[0].pending.event.name = EventName{value: "account.closed"}
		}
		if err := lifecycle.Acknowledge(changes, prepared, messages); !errors.Is(
			err,
			ErrPersistenceMismatch,
		) || !lifecycle.poisoned {
			t.Fatalf(
				"Acknowledge(persistence mismatch %d) error = %v, poisoned = %t",
				index,
				err,
				lifecycle.poisoned,
			)
		}
	}
}

func TestHistoryValidationChecksEveryCoordinateIndependently(t *testing.T) {
	t.Parallel()

	event := invariantDecodedEvent()
	valid := HistoricalEvent{
		sourceVersion: 1,
		segmentCount:  1,
		event:         event,
	}
	firstCases := [][]HistoricalEvent{
		{{sourceVersion: 1, segmentCount: 1}},
		{{sourceVersion: 2, segmentCount: 1, event: event}},
		{{sourceVersion: 1, segmentIndex: 1, segmentCount: 1, event: event}},
		{{sourceVersion: 1, event: event}},
		{{sourceVersion: 1, segmentCount: MaxUpcastSegments + 1, event: event}},
	}
	for index, history := range firstCases {
		if err := validateHistory(0, history); !errors.Is(err, ErrCorruptHistory) {
			t.Fatalf("validateHistory(first mismatch %d) error = %v", index, err)
		}
	}

	maximumSplit := make([]HistoricalEvent, MaxUpcastSegments)
	for index := range maximumSplit {
		maximumSplit[index] = HistoricalEvent{
			sourceVersion: 1,
			segmentIndex:  uint32(index),
			segmentCount:  MaxUpcastSegments,
			event:         event,
		}
	}
	if err := validateHistory(0, maximumSplit); err != nil {
		t.Fatalf("validateHistory(maximum split) error = %v", err)
	}

	incompleteLaterSplit := []HistoricalEvent{
		valid,
		{sourceVersion: 2, segmentCount: 2, event: event},
	}
	if err := validateHistory(0, incompleteLaterSplit); !errors.Is(
		err,
		ErrCorruptHistory,
	) {
		t.Fatalf("validateHistory(incomplete later split) error = %v", err)
	}

	second := HistoricalEvent{
		sourceVersion: 1,
		segmentIndex:  1,
		segmentCount:  2,
		event:         event,
	}
	first := second
	first.segmentIndex = 0
	innerCases := []HistoricalEvent{second, second, second, second}
	innerCases[0].event = DecodedEvent{}
	innerCases[1].sourceVersion = 2
	innerCases[2].segmentIndex = 0
	innerCases[3].segmentCount = 1
	for index, candidate := range innerCases {
		if err := validateHistory(0, []HistoricalEvent{first, candidate}); !errors.Is(
			err,
			ErrCorruptHistory,
		) {
			t.Fatalf("validateHistory(inner mismatch %d) error = %v", index, err)
		}
	}
}

func lifecycleChangeSetFixture() (*Lifecycle, ChangeSet) {
	event := invariantDecodedEvent()
	lifecycle := &Lifecycle{
		committed:  2,
		pending:    []DecodedEvent{event},
		generation: 7,
	}
	changes := ChangeSet{
		owner:      lifecycle,
		generation: lifecycle.generation,
		base:       lifecycle.committed,
		events:     []DecodedEvent{event},
	}

	return lifecycle, changes
}

func acknowledgementFixture() (*Lifecycle, ChangeSet, []PendingMessage, []Message) {
	event := invariantDecodedEvent()
	lifecycle := &Lifecycle{
		committed:  2,
		pending:    []DecodedEvent{event},
		generation: 7,
	}
	changes := ChangeSet{
		owner:      lifecycle,
		generation: lifecycle.generation,
		base:       lifecycle.committed,
		events:     []DecodedEvent{event},
	}
	pending := invariantPendingMessage()
	message := Message{pending: pending, streamVersion: 3}

	return lifecycle, changes, []PendingMessage{pending}, []Message{message}
}
