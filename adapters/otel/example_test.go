package otel_test

import (
	"context"
	"errors"
	"fmt"
	"log"

	eventotel "github.com/faustbrian/go-event-sourcing/adapters/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type exampleRuntime struct {
	tracer     trace.TracerProvider
	meter      metric.MeterProvider
	propagator propagation.TextMapPropagator
}

func (runtime exampleRuntime) TracerProvider() trace.TracerProvider {
	return runtime.tracer
}

func (runtime exampleRuntime) MeterProvider() metric.MeterProvider {
	return runtime.meter
}

func (runtime exampleRuntime) Propagator() propagation.TextMapPropagator {
	return runtime.propagator
}

func ExampleNew() {
	tracerProvider := sdktrace.NewTracerProvider()
	meterProvider := sdkmetric.NewMeterProvider()
	runtime := exampleRuntime{
		tracer:     tracerProvider,
		meter:      meterProvider,
		propagator: propagation.TraceContext{},
	}

	instrumentation, err := eventotel.New(runtime)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("instrumentation ready", instrumentation != nil)

	if err := errors.Join(
		meterProvider.Shutdown(context.Background()),
		tracerProvider.Shutdown(context.Background()),
	); err != nil {
		log.Fatal(err)
	}

	// Output: instrumentation ready true
}
