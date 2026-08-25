package projection

import (
	"errors"
	"testing"

	eventsourcing "github.com/faustbrian/go-event-sourcing"
)

func TestPoisonedDeliveryRequiresDeliveryAndCauseIndependently(t *testing.T) {
	t.Parallel()

	delivery, err := eventsourcing.NewDelivery(
		internalProjectionMessage(t, 1),
		eventsourcing.DeliveryReplay,
	)
	if err != nil {
		t.Fatal(err)
	}
	cause := errors.New("handler failure")
	if !(PoisonedDelivery{delivery: delivery}).IsZero() ||
		!(PoisonedDelivery{cause: cause}).IsZero() ||
		(PoisonedDelivery{delivery: delivery, cause: cause}).IsZero() {
		t.Fatal("PoisonedDelivery zero-state components are not independent")
	}
}

func TestNewRunnerAcceptsExactBatchBound(t *testing.T) {
	t.Parallel()

	config := internalRunnerConfig()
	config.BatchSize = eventsourcing.MaxReadMessages
	runner, err := NewRunner(config)
	if err != nil || runner == nil || runner.batchSize != eventsourcing.MaxReadMessages {
		t.Fatalf("NewRunner(exact batch bound) = %#v, %v", runner, err)
	}
}
