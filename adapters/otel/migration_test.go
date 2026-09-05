package otel_test

import gotelemetry "github.com/faustbrian/go-event-sourcing/adapters/otel"

// Compile representative selectors with the compatibility alias documented
// for consumers migrating from adapters/gotelemetry.
var (
	_ = gotelemetry.New
	_ = gotelemetry.WrapProcessManager[struct{}]
)
