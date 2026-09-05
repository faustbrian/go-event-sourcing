package kafka

import "testing"

func TestSuccessorExposesKafkaAdapterContract(t *testing.T) {
	limits := DefaultRecordLimits()
	if err := limits.Validate(); err != nil {
		t.Fatalf("DefaultRecordLimits() error = %v", err)
	}
}
