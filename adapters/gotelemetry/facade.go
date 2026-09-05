// Package gotelemetry preserves the original OpenTelemetry adapter import path.
//
// Deprecated: use github.com/faustbrian/go-event-sourcing/adapters/otel.
package gotelemetry

import successor "github.com/faustbrian/go-event-sourcing/adapters/otel"

var (
	ErrRuntimeRequired                   = successor.ErrRuntimeRequired
	ErrDispatcherRequired                = successor.ErrDispatcherRequired
	ErrConsumerRequired                  = successor.ErrConsumerRequired
	ErrContextRequired                   = successor.ErrContextRequired
	ErrInstrumentCreation                = successor.ErrInstrumentCreation
	ErrKafkaPublisherRequired            = successor.ErrKafkaPublisherRequired
	ErrKafkaHandlerRequired              = successor.ErrKafkaHandlerRequired
	ErrInvalidKafkaPropagation           = successor.ErrInvalidKafkaPropagation
	ErrKafkaPropagationRejected          = successor.ErrKafkaPropagationRejected
	ErrProcessManagerNameInvalid         = successor.ErrProcessManagerNameInvalid
	ErrProcessManagerRequired            = successor.ErrProcessManagerRequired
	ErrProjectionNameInvalid             = successor.ErrProjectionNameInvalid
	ErrProjectionRunnerRequired          = successor.ErrProjectionRunnerRequired
	ErrProjectionLagInvalid              = successor.ErrProjectionLagInvalid
	ErrPayloadCodecRequired              = successor.ErrPayloadCodecRequired
	ErrUpcasterRequired                  = successor.ErrUpcasterRequired
	ErrEventStoreRequired                = successor.ErrEventStoreRequired
	ErrGlobalReaderRequired              = successor.ErrGlobalReaderRequired
	ErrMessageIteratorRequired           = successor.ErrMessageIteratorRequired
	ErrProjectionCheckpointStoreRequired = successor.ErrProjectionCheckpointStoreRequired
	ErrProjectionControllerRequired      = successor.ErrProjectionControllerRequired
	ErrProjectionHandlerRequired         = successor.ErrProjectionHandlerRequired
	ErrSnapshotStoreRequired             = successor.ErrSnapshotStoreRequired
)

type (
	InstrumentError        = successor.InstrumentError
	Instrumentation        = successor.Instrumentation
	KafkaPropagationConfig = successor.KafkaPropagationConfig
	KafkaPropagationError  = successor.KafkaPropagationError
	KafkaPublisher         = successor.KafkaPublisher
	ProjectionController   = successor.ProjectionController
	ProjectionRunner       = successor.ProjectionRunner
	Runtime                = successor.Runtime
)

type ProcessManager[Command any] = successor.ProcessManager[Command]

func New(runtime Runtime) (*Instrumentation, error) { return successor.New(runtime) }

func WrapProcessManager[Command any](instrumentation *Instrumentation, name string, next ProcessManager[Command]) (ProcessManager[Command], error) {
	return successor.WrapProcessManager(instrumentation, name, next)
}
