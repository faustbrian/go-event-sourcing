package otel

import (
	"errors"
	"testing"
)

func TestSuccessorExposesOpenTelemetryAdapterContract(t *testing.T) {
	_, err := New(nil)
	if !errors.Is(err, ErrRuntimeRequired) {
		t.Fatalf("New(nil) error = %v", err)
	}
}
