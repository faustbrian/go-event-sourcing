package eventsourcing_test

import (
	"context"
	"errors"
	"testing"

	eventsourcing "github.com/faustbrian/go-event-sourcing"
	"github.com/faustbrian/go-event-sourcing/memory"
)

func TestAggregateRepositoryAcceptsExactConfiguredAndAppendBounds(t *testing.T) {
	t.Parallel()

	config := completeRepositoryConfig(t, memory.NewStore())
	config.ReadBatchSize = eventsourcing.MaxReadMessages
	repository := repositoryFromConfig(t, config)
	account := &repositoryAccount{id: "account-maximum"}
	for range eventsourcing.MaxAppendMessages {
		if err := account.lifecycle.Record(repositoryOpenedEvent(t), account.apply); err != nil {
			t.Fatal(err)
		}
	}
	plan, err := repository.PrepareSave(context.Background(), account)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.PreparedMessages()) != eventsourcing.MaxAppendMessages {
		t.Fatalf(
			"PrepareSave() length = %d, want %d",
			len(plan.PreparedMessages()),
			eventsourcing.MaxAppendMessages,
		)
	}
}

func TestAggregateRepositoryContinuesAfterAnExactlyFullPage(t *testing.T) {
	t.Parallel()

	stream, err := eventsourcing.NewStreamID("account", "account-42")
	if err != nil {
		t.Fatal(err)
	}
	reads := 0
	store := &repositoryStoreFuncs{
		read: func(
			_ context.Context,
			_ eventsourcing.StreamID,
			options eventsourcing.ReadStreamOptions,
		) (eventsourcing.MessageIterator, error) {
			reads++
			switch reads {
			case 1:
				if options.FromVersion() != 1 || options.Limit() != 2 {
					t.Fatalf("first ReadStream() options = %#v", options)
				}

				return &repositorySliceIterator{messages: []eventsourcing.Message{
					repositoryMessageAt(t, stream, 1),
					repositoryMessageAt(t, stream, 2),
				}}, nil
			case 2:
				if options.FromVersion() != 3 {
					t.Fatalf("second ReadStream() options = %#v", options)
				}

				return &repositorySliceIterator{}, nil
			default:
				return nil, errors.New("unexpected read after empty page")
			}
		},
	}
	config := completeRepositoryConfig(t, store)
	config.ReadBatchSize = 2
	repository := repositoryFromConfig(t, config)
	aggregate, err := repository.Load(context.Background(), "account-42")
	if err != nil || aggregate == nil || reads != 2 ||
		aggregate.lifecycle.CommittedVersion() != 2 {
		t.Fatalf("Load() = %#v, %v, reads = %d", aggregate, err, reads)
	}
}

func TestAggregateRepositoryStopsAtMaximumCommittedVersion(t *testing.T) {
	t.Parallel()

	const base = ^uint64(0) - 2
	stream, err := eventsourcing.NewStreamID("account", "account-42")
	if err != nil {
		t.Fatal(err)
	}
	reads := 0
	store := &repositoryStoreFuncs{
		read: func(
			_ context.Context,
			_ eventsourcing.StreamID,
			options eventsourcing.ReadStreamOptions,
		) (eventsourcing.MessageIterator, error) {
			reads++
			switch reads {
			case 1:
				if options.FromVersion() != base || options.ToVersion() != base {
					t.Fatalf("snapshot verification options = %#v", options)
				}

				return &repositorySliceIterator{messages: []eventsourcing.Message{
					repositoryMessageAt(t, stream, base),
				}}, nil
			case 2:
				if options.FromVersion() != base+1 || options.Limit() != 2 {
					t.Fatalf("history options = %#v", options)
				}

				return &repositorySliceIterator{messages: []eventsourcing.Message{
					repositoryMessageAt(t, stream, base+1),
					repositoryMessageAt(t, stream, ^uint64(0)),
				}}, nil
			default:
				return nil, errors.New("read advanced past maximum version")
			}
		},
	}
	config := completeRepositoryConfig(t, store)
	config.ReadBatchSize = 2
	repository := repositoryFromConfig(t, config)
	aggregate, err := repository.Restore(
		context.Background(),
		"account-42",
		base,
		func() (*repositoryAccount, error) {
			return &repositoryAccount{id: "account-42", owner: "Ada"}, nil
		},
	)
	if err != nil || aggregate == nil || reads != 2 ||
		aggregate.lifecycle.CommittedVersion() != ^uint64(0) {
		t.Fatalf("Restore() = %#v, %v, reads = %d", aggregate, err, reads)
	}
}

func TestAggregateRepositoryStopsReadingAtTheFirstCorruptBoundary(t *testing.T) {
	t.Parallel()

	stream, err := eventsourcing.NewStreamID("account", "account-42")
	if err != nil {
		t.Fatal(err)
	}
	tests := map[string]struct {
		messages []eventsourcing.Message
		batch    uint32
		wantNext int
	}{
		"page overrun": {
			messages: []eventsourcing.Message{
				repositoryMessageAt(t, stream, 1),
				repositoryMessageAt(t, stream, 2),
			},
			batch:    1,
			wantNext: 2,
		},
		"wrong first version": {
			messages: []eventsourcing.Message{
				repositoryMessageAt(t, stream, 2),
				repositoryMessageAt(t, stream, 1),
			},
			batch:    2,
			wantNext: 1,
		},
	}
	for name, testCase := range tests {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			iterator := &countingRepositoryIterator{messages: testCase.messages}
			config := completeRepositoryConfig(t, repositoryReadStore(iterator))
			config.ReadBatchSize = testCase.batch
			repository := repositoryFromConfig(t, config)
			aggregate, err := repository.Load(context.Background(), "account-42")
			if aggregate != nil || !errors.Is(err, eventsourcing.ErrCorruptHistory) ||
				iterator.nextCalls != testCase.wantNext {
				t.Fatalf(
					"Load() = %#v, %v, Next calls = %d, want %d",
					aggregate,
					err,
					iterator.nextCalls,
					testCase.wantNext,
				)
			}
		})
	}
}

func TestAggregateRepositoryStopsAfterTheFirstDecodeFailure(t *testing.T) {
	t.Parallel()

	stream, err := eventsourcing.NewStreamID("account", "account-42")
	if err != nil {
		t.Fatal(err)
	}
	iterator := &countingRepositoryIterator{messages: []eventsourcing.Message{
		repositoryMessageAt(t, stream, 1),
		repositoryMessageAt(t, stream, 2),
	}}
	decodeFailure := errors.New("decode failure")
	decodeCalls := 0
	config := completeRepositoryConfig(t, repositoryReadStore(iterator))
	config.ReadBatchSize = 2
	config.Codec = repositoryCodecFuncs{
		decode: func(eventsourcing.EncodedEvent) (eventsourcing.DecodedEvent, error) {
			decodeCalls++

			return eventsourcing.DecodedEvent{}, decodeFailure
		},
	}
	repository := repositoryFromConfig(t, config)
	aggregate, err := repository.Load(context.Background(), "account-42")
	if aggregate != nil || !errors.Is(err, decodeFailure) ||
		decodeCalls != 1 || iterator.nextCalls != 1 {
		t.Fatalf(
			"Load() = %#v, %v, decode calls = %d, Next calls = %d",
			aggregate,
			err,
			decodeCalls,
			iterator.nextCalls,
		)
	}
}

type countingRepositoryIterator struct {
	messages  []eventsourcing.Message
	index     int
	nextCalls int
	current   eventsourcing.Message
}

func (iterator *countingRepositoryIterator) Next(context.Context) bool {
	iterator.nextCalls++
	if iterator.index >= len(iterator.messages) {
		return false
	}
	iterator.current = iterator.messages[iterator.index]
	iterator.index++

	return true
}

func (iterator *countingRepositoryIterator) Message() eventsourcing.Message {
	return iterator.current
}

func (*countingRepositoryIterator) Err() error {
	return nil
}

func (*countingRepositoryIterator) Close() error {
	return nil
}
