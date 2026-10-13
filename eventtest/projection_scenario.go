package eventtest

import (
	"context"
	"errors"
	"fmt"

	eventsourcing "github.com/faustbrian/go-event-sourcing/v2"
	"github.com/faustbrian/go-event-sourcing/v2/projection"
)

// ProjectionRunner executes one bounded projection batch.
type ProjectionRunner interface {
	// RunBatch returns handled and checkpointed progress separately, including
	// partial progress on error.
	RunBatch(context.Context) (projection.BatchResult, error)
}

// ExpectedProjectionBatch describes the observable progress from one runner
// batch.
type ExpectedProjectionBatch struct {
	// Scanned is the expected number of persisted messages examined.
	Scanned uint32
	// Handled is the expected successful handler count, even if checkpointing
	// later failed.
	Handled uint32
	// Filtered is the expected deliberately omitted handler count.
	Filtered uint32
	// Skipped is the expected explicitly durably skipped poison-message count.
	Skipped uint32
	// Checkpointed is the expected number followed by successful checkpoint saves,
	// including filtered records.
	Checkpointed uint32
	// Checkpoint is the expected last durable global position, not just the last
	// handled position.
	Checkpoint eventsourcing.GlobalPosition
}

// ProjectionScenario describes one bounded runner call and optional
// application-owned read-model predicate. WantError uses errors.Is semantics.
type ProjectionScenario struct {
	// Runner is the required borrowed boundary called once using the supplied
	// context.
	Runner ProjectionRunner
	// Expected contains all counters and the durable checkpoint compared even when
	// an error is expected.
	Expected ExpectedProjectionBatch
	// WantError selects errors.Is comparison; nil expects success, not a string
	// comparison.
	WantError error
	// State is an optional read-model predicate invoked after progress/error
	// comparisons; nil omits that assertion.
	State func() bool
}

// CheckProjectionScenario executes one bounded projection batch and compares
// durable progress plus application state without formatting event or model
// data.
func CheckProjectionScenario(
	ctx context.Context,
	scenario ProjectionScenario,
) error {
	if ctx == nil || scenario.Runner == nil {
		return eventsourcing.ErrInvalidArgument
	}

	result, err := scenario.Runner.RunBatch(ctx)
	if scenario.WantError == nil {
		if err != nil {
			return err
		}
	} else if !errors.Is(err, scenario.WantError) {
		return fmt.Errorf(
			"%w: projection error category differs",
			ErrConformance,
		)
	}
	if result.Scanned() != scenario.Expected.Scanned {
		return projectionMismatch("scanned count")
	}
	if result.Handled() != scenario.Expected.Handled {
		return projectionMismatch("handled count")
	}
	if result.Filtered() != scenario.Expected.Filtered {
		return projectionMismatch("filtered count")
	}
	if result.Skipped() != scenario.Expected.Skipped {
		return projectionMismatch("skipped count")
	}
	if result.Checkpointed() != scenario.Expected.Checkpointed {
		return projectionMismatch("checkpointed count")
	}
	if result.Checkpoint() != scenario.Expected.Checkpoint {
		return projectionMismatch("checkpoint")
	}
	if scenario.State != nil && !scenario.State() {
		return projectionMismatch("application state")
	}

	return nil
}

func projectionMismatch(field string) error {
	return fmt.Errorf("%w: projection %s differs", ErrConformance, field)
}

var _ ProjectionRunner = (*projection.Runner)(nil)
