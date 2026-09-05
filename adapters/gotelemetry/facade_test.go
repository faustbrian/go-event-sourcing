package gotelemetry_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	eventsourcing "github.com/faustbrian/go-event-sourcing"
	legacy "github.com/faustbrian/go-event-sourcing/adapters/gotelemetry"
	successor "github.com/faustbrian/go-event-sourcing/adapters/otel"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

func ExampleNew() {
	instrumentation, err := legacy.New(testRuntime{
		tracer: tracenoop.NewTracerProvider(),
		meter:  metricnoop.NewMeterProvider(),
	})
	fmt.Println(instrumentation != nil, err)

	// Output: true <nil>
}

func TestFacadePreservesCategoriesAndConstruction(t *testing.T) {
	if !errors.Is(legacy.ErrRuntimeRequired, successor.ErrRuntimeRequired) {
		t.Fatal("legacy error category differs from successor")
	}
	_, err := legacy.New(nil)
	if !errors.Is(err, legacy.ErrRuntimeRequired) {
		t.Fatalf("New(nil) error = %v", err)
	}
	_, err = legacy.WrapProcessManager[string](nil, "manager", nil)
	if !errors.Is(err, legacy.ErrRuntimeRequired) {
		t.Fatalf("WrapProcessManager(nil) error = %v", err)
	}
}

func TestFacadePreservesInstrumentationScope(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	instrumentation, err := legacy.New(testRuntime{
		tracer: sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder)),
		meter:  sdkmetric.NewMeterProvider(),
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	consumer, err := instrumentation.WrapConsumer(func(context.Context, eventsourcing.Delivery) error {
		return nil
	})
	if err != nil {
		t.Fatalf("WrapConsumer() error = %v", err)
	}
	if err := consumer(context.Background(), testDelivery(t)); err != nil {
		t.Fatalf("consumer() error = %v", err)
	}
	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("ended spans = %d, want 1", len(spans))
	}
	scope := spans[0].InstrumentationScope()
	if scope.Name != "github.com/faustbrian/go-event-sourcing" || scope.Version != "" || scope.SchemaURL != "" {
		t.Fatalf("instrumentation scope = %#v", scope)
	}
}

type testRuntime struct {
	tracer trace.TracerProvider
	meter  metric.MeterProvider
}

func (runtime testRuntime) TracerProvider() trace.TracerProvider { return runtime.tracer }
func (runtime testRuntime) MeterProvider() metric.MeterProvider  { return runtime.meter }
func (testRuntime) Propagator() propagation.TextMapPropagator    { return propagation.TraceContext{} }

func testDelivery(t *testing.T) eventsourcing.Delivery {
	t.Helper()
	stream, err := eventsourcing.NewStreamID("account", "42")
	if err != nil {
		t.Fatal(err)
	}
	event, err := eventsourcing.NewEncodedEvent(eventsourcing.EncodedEventInput{
		Name: "account.changed", Version: 1, ContentType: "application/json", Payload: []byte(`{}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	pending, err := eventsourcing.NewPendingMessage(eventsourcing.PendingMessageInput{
		ID: "message-1", Stream: stream, Event: event, RecordedAt: time.Unix(1, 0),
	})
	if err != nil {
		t.Fatal(err)
	}
	message, err := eventsourcing.NewMessage(eventsourcing.MessageInput{Pending: pending, StreamVersion: 1})
	if err != nil {
		t.Fatal(err)
	}
	delivery, err := eventsourcing.NewDelivery(message, eventsourcing.DeliveryLive)
	if err != nil {
		t.Fatal(err)
	}
	return delivery
}
