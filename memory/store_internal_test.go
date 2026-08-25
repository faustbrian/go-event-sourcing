package memory

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	eventsourcing "github.com/faustbrian/go-event-sourcing"
)

func TestStoreChecksCancellationAfterWaitingForWriteOwnership(t *testing.T) {
	t.Parallel()

	store := NewStore()
	stream := internalStream(t)
	ctx := &cancelAfterChecks{allowed: 1}

	if _, err := store.Append(
		ctx,
		stream,
		eventsourcing.ExpectNewStream(),
		[]eventsourcing.PendingMessage{internalPending(t, stream)},
	); !errors.Is(err, context.Canceled) {
		t.Fatalf("Append() error = %v, want context.Canceled", err)
	}
}

func TestStoreChecksCancellationWhilePreparingAtomicBatch(t *testing.T) {
	t.Parallel()

	store := NewStore()
	stream := internalStream(t)
	ctx := &cancelAfterChecks{allowed: 2}

	if _, err := store.Append(
		ctx,
		stream,
		eventsourcing.ExpectNewStream(),
		[]eventsourcing.PendingMessage{internalPending(t, stream)},
	); !errors.Is(err, context.Canceled) {
		t.Fatalf("Append() error = %v, want context.Canceled", err)
	}
	if len(store.streams) != 0 || len(store.messageIDs) != 0 {
		t.Fatal("cancelled batch mutated store")
	}
}

func TestStoreRejectsGlobalPositionOverflow(t *testing.T) {
	t.Parallel()

	store := NewStore()
	store.globalPosition = eventsourcing.GlobalPosition(^uint64(0))
	stream := internalStream(t)

	if _, err := store.Append(
		context.Background(),
		stream,
		eventsourcing.ExpectNewStream(),
		[]eventsourcing.PendingMessage{internalPending(t, stream)},
	); !errors.Is(err, eventsourcing.ErrVersionOverflow) {
		t.Fatalf("Append() error = %v, want ErrVersionOverflow", err)
	}
}

func TestStoreRequiresBothInternalIndexes(t *testing.T) {
	t.Parallel()

	stream := internalStream(t)
	pending := []eventsourcing.PendingMessage{internalPending(t, stream)}
	options, err := eventsourcing.NewReadStreamOptions(eventsourcing.ReadStreamOptionsInput{
		FromVersion: 1,
		Limit:       1,
	})
	if err != nil {
		t.Fatal(err)
	}
	globalOptions, err := eventsourcing.NewReadGlobalOptions(
		eventsourcing.ReadGlobalOptionsInput{FromPosition: 1, Limit: 1},
	)
	if err != nil {
		t.Fatal(err)
	}

	stores := []*Store{
		{streams: make(map[eventsourcing.StreamID][]eventsourcing.Message)},
		{messageIDs: make(map[string]struct{})},
	}
	for index, store := range stores {
		if _, appendErr := store.Append(
			context.Background(),
			stream,
			eventsourcing.ExpectNewStream(),
			pending,
		); !errors.Is(appendErr, eventsourcing.ErrInvalidArgument) {
			t.Fatalf("store %d Append() error = %v", index, appendErr)
		}
		if _, readErr := store.ReadStream(
			context.Background(),
			stream,
			options,
		); !errors.Is(readErr, eventsourcing.ErrInvalidArgument) {
			t.Fatalf("store %d ReadStream() error = %v", index, readErr)
		}
		if _, readErr := store.ReadGlobal(
			context.Background(),
			globalOptions,
		); !errors.Is(readErr, eventsourcing.ErrInvalidArgument) {
			t.Fatalf("store %d ReadGlobal() error = %v", index, readErr)
		}
	}
}

func TestStoreAcceptsExactAppendAndGlobalPositionBounds(t *testing.T) {
	t.Parallel()

	stream := internalStream(t)
	messages := make([]eventsourcing.PendingMessage, eventsourcing.MaxAppendMessages)
	for index := range messages {
		messages[index] = internalPendingID(t, stream, fmt.Sprintf("message-%d", index))
	}
	store := NewStore()
	stored, err := store.Append(
		context.Background(),
		stream,
		eventsourcing.ExpectNewStream(),
		messages,
	)
	if err != nil || len(stored) != eventsourcing.MaxAppendMessages ||
		store.globalPosition != eventsourcing.GlobalPosition(eventsourcing.MaxAppendMessages) {
		t.Fatalf("Append(maximum batch) = %d, %v, position = %d", len(stored), err, store.globalPosition)
	}

	boundary := NewStore()
	boundary.globalPosition = eventsourcing.GlobalPosition(^uint64(0) - 1)
	stored, err = boundary.Append(
		context.Background(),
		stream,
		eventsourcing.ExpectNewStream(),
		[]eventsourcing.PendingMessage{internalPending(t, stream)},
	)
	if err != nil || len(stored) != 1 {
		t.Fatalf("Append(exact global bound) = %#v, %v", stored, err)
	}
	position, assigned := stored[0].GlobalPosition()
	if !assigned || position != eventsourcing.GlobalPosition(^uint64(0)) ||
		boundary.globalPosition != eventsourcing.GlobalPosition(^uint64(0)) {
		t.Fatalf("Append(exact global bound) = %#v, %v, position = %d", stored, err, boundary.globalPosition)
	}
}

func TestStoreInternalReadRangesPreserveExactBounds(t *testing.T) {
	t.Parallel()

	store := NewStore()
	stream := internalStream(t)
	messages := []eventsourcing.PendingMessage{
		internalPendingID(t, stream, "message-1"),
		internalPendingID(t, stream, "message-2"),
		internalPendingID(t, stream, "message-3"),
	}
	if _, err := store.Append(
		context.Background(),
		stream,
		eventsourcing.ExpectNewStream(),
		messages,
	); err != nil {
		t.Fatal(err)
	}

	streamCases := []struct {
		input eventsourcing.ReadStreamOptionsInput
		want  int
	}{
		{input: eventsourcing.ReadStreamOptionsInput{FromVersion: 1, Limit: 3}, want: 3},
		{input: eventsourcing.ReadStreamOptionsInput{FromVersion: 1, ToVersion: 2, Limit: 3}, want: 2},
		{input: eventsourcing.ReadStreamOptionsInput{FromVersion: 1, Limit: 2}, want: 2},
		{input: eventsourcing.ReadStreamOptionsInput{FromVersion: 4, Limit: 1}},
	}
	for index, testCase := range streamCases {
		options, err := eventsourcing.NewReadStreamOptions(testCase.input)
		if err != nil {
			t.Fatal(err)
		}
		iterator, err := store.ReadStream(context.Background(), stream, options)
		if err != nil {
			t.Fatal(err)
		}
		if got := internalIteratorCount(t, iterator); got != testCase.want {
			t.Fatalf("stream range %d count = %d, want %d", index, got, testCase.want)
		}
	}

	globalCases := []struct {
		input eventsourcing.ReadGlobalOptionsInput
		want  int
	}{
		{input: eventsourcing.ReadGlobalOptionsInput{FromPosition: 1, Limit: 3}, want: 3},
		{input: eventsourcing.ReadGlobalOptionsInput{FromPosition: 1, ToPosition: 2, Limit: 3}, want: 2},
		{input: eventsourcing.ReadGlobalOptionsInput{FromPosition: 1, Limit: 2}, want: 2},
		{input: eventsourcing.ReadGlobalOptionsInput{FromPosition: 4, Limit: 1}},
	}
	for index, testCase := range globalCases {
		options, err := eventsourcing.NewReadGlobalOptions(testCase.input)
		if err != nil {
			t.Fatal(err)
		}
		iterator, err := store.ReadGlobal(context.Background(), options)
		if err != nil {
			t.Fatal(err)
		}
		if got := internalIteratorCount(t, iterator); got != testCase.want {
			t.Fatalf("global range %d count = %d, want %d", index, got, testCase.want)
		}
	}
}

type cancelAfterChecks struct {
	allowed int
	checks  int
}

func (*cancelAfterChecks) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (*cancelAfterChecks) Done() <-chan struct{} {
	return nil
}

func (ctx *cancelAfterChecks) Err() error {
	ctx.checks++
	if ctx.checks > ctx.allowed {
		return context.Canceled
	}

	return nil
}

func (*cancelAfterChecks) Value(any) any {
	return nil
}

func internalStream(t *testing.T) eventsourcing.StreamID {
	t.Helper()

	stream, err := eventsourcing.NewStreamID("bank.account", "account-42")
	if err != nil {
		t.Fatal(err)
	}

	return stream
}

func internalPending(
	t *testing.T,
	stream eventsourcing.StreamID,
) eventsourcing.PendingMessage {
	return internalPendingID(t, stream, "message-1")
}

func internalPendingID(
	t *testing.T,
	stream eventsourcing.StreamID,
	id string,
) eventsourcing.PendingMessage {
	t.Helper()

	event, err := eventsourcing.NewEncodedEvent(eventsourcing.EncodedEventInput{
		Name:        "account.opened",
		Version:     1,
		ContentType: eventsourcing.JSONContentType,
		Payload:     []byte("{}"),
	})
	if err != nil {
		t.Fatal(err)
	}
	pending, err := eventsourcing.NewPendingMessage(eventsourcing.PendingMessageInput{
		ID:         id,
		Stream:     stream,
		Event:      event,
		RecordedAt: time.Date(2026, time.July, 25, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}

	return pending
}

func internalIteratorCount(t *testing.T, iterator eventsourcing.MessageIterator) int {
	t.Helper()

	count := 0
	for iterator.Next(context.Background()) {
		count++
	}
	if err := iterator.Err(); err != nil {
		t.Fatal(err)
	}
	if err := iterator.Close(); err != nil {
		t.Fatal(err)
	}

	return count
}
