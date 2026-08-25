package processmanager

import (
	"context"
	"fmt"
	"testing"

	eventsourcing "github.com/faustbrian/go-event-sourcing"
)

func TestNewAcceptsExactConfigurationBounds(t *testing.T) {
	t.Parallel()

	names := make([]eventsourcing.EventName, MaxAcceptedEventNames)
	for index := range names {
		name, err := eventsourcing.NewEventName(fmt.Sprintf("event.type-%d", index))
		if err != nil {
			t.Fatal(err)
		}
		names[index] = name
	}
	manager, err := New(Config[struct{}]{
		Name:        "bounded-process",
		EventNames:  names,
		MaxCommands: MaxPlannedCommands,
		Planner: func(context.Context, eventsourcing.Delivery) ([]struct{}, error) {
			return nil, nil
		},
	})
	if err != nil || manager == nil || len(manager.eventNames) != MaxAcceptedEventNames ||
		manager.maxCommands != MaxPlannedCommands {
		t.Fatalf("New(exact bounds) = %#v, %v", manager, err)
	}
}
