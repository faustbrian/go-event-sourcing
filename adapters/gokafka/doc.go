// Package gokafka preserves the original Kafka adapter import path.
//
// This independently published compatibility module retains core-v1 types and
// its original import identity. Its synchronous codec/dispatch/handling contract
// uses the same package-defined wire mapping as the canonical adapter; broker
// configuration, offset settlement, retries and publisher lifetime remain caller-owned.
//
// Deprecated: use github.com/faustbrian/go-event-sourcing/adapters/kafka/v2.
package gokafka
