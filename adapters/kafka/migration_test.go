package kafka_test

import gokafka "github.com/faustbrian/go-event-sourcing/adapters/kafka"

// Compile representative selectors with the compatibility alias documented
// for consumers migrating from adapters/gokafka.
var (
	_ = gokafka.NewRecordCodec
	_ = gokafka.NewDispatcher
	_ = gokafka.NewRecordHandler
)
