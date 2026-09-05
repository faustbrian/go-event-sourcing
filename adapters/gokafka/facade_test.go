package gokafka_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	eventsourcing "github.com/faustbrian/go-event-sourcing"
	legacy "github.com/faustbrian/go-event-sourcing/adapters/gokafka"
	successor "github.com/faustbrian/go-event-sourcing/adapters/kafka"
	broker "github.com/faustbrian/go-kafka"
)

func ExampleNewRecordCodec() {
	codec, err := legacy.NewRecordCodec(legacy.RecordCodecConfig{
		Resolver:      legacy.FixedTopic("accounts.events.v1"),
		AllowedTopics: []string{"accounts.events.v1"},
	})
	fmt.Println(codec != nil, err)

	// Output: true <nil>
}

func TestFacadePreservesCategoriesAndConstruction(t *testing.T) {
	if !errors.Is(legacy.ErrInvalidConfig, successor.ErrInvalidConfig) {
		t.Fatal("legacy error category differs from successor")
	}
	codec, err := legacy.NewRecordCodec(legacy.RecordCodecConfig{
		Resolver:      legacy.FixedTopic("events.v1"),
		AllowedTopics: []string{"events.v1"},
	})
	if err != nil || codec == nil {
		t.Fatalf("NewRecordCodec() = %v, %v", codec, err)
	}

	publisher := publisherFunc(func(context.Context, broker.Message) error { return nil })
	dispatcher, err := legacy.NewDispatcher(
		publisher,
		codec,
		legacy.AllowReplay(),
		legacy.ContinueOnPublishError(),
	)
	if err != nil || dispatcher == nil {
		t.Fatalf("NewDispatcher() = %v, %v", dispatcher, err)
	}

	policy := legacy.FailurePolicyFunc(func(
		context.Context,
		broker.ConsumedMessage,
		error,
	) (legacy.FailureDisposition, error) {
		return legacy.FailureRetry, nil
	})
	handler, err := legacy.NewRecordHandler(
		codec,
		legacy.DeliveryConsumerFunc(func(context.Context, eventsourcing.Delivery) error { return nil }),
		legacy.AllowReplayHandling(),
		legacy.WithFailurePolicy(policy),
	)
	if err != nil || handler == nil {
		t.Fatalf("NewRecordHandler() = %v, %v", handler, err)
	}

	deadLetter, err := legacy.NewDeadLetterPolicy(publisher, legacy.DeadLetterPolicyConfig{
		Topic:  "events.dead-letter.v1",
		Limits: legacy.DefaultRecordLimits(),
	})
	if err != nil || deadLetter == nil {
		t.Fatalf("NewDeadLetterPolicy() = %v, %v", deadLetter, err)
	}
}

type publisherFunc func(context.Context, broker.Message) error

func (publish publisherFunc) Publish(ctx context.Context, message broker.Message) error {
	return publish(ctx, message)
}
