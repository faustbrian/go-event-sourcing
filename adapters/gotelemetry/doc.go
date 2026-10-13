// Package gotelemetry preserves the original OpenTelemetry adapter import path.
//
// This independently published compatibility module retains core-v1 types,
// imports and instrumentation scope. Providers/instruments/exporter shutdown remain
// caller-owned. Provider panic isolation does not impose a deadline; downstream
// operation errors, panics and partial results pass through without diagnostics
// being recorded by the adapter.
//
// Deprecated: use github.com/faustbrian/go-event-sourcing/adapters/otel/v2.
package gotelemetry
