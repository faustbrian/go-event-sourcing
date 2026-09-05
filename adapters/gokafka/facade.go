// Package gokafka preserves the original Kafka adapter import path.
//
// Deprecated: use github.com/faustbrian/go-event-sourcing/adapters/kafka.
package gokafka

import (
	successor "github.com/faustbrian/go-event-sourcing/adapters/kafka"
	broker "github.com/faustbrian/go-kafka"
)

const (
	HeaderWireVersion               = successor.HeaderWireVersion
	HeaderMessageID                 = successor.HeaderMessageID
	HeaderAggregateType             = successor.HeaderAggregateType
	HeaderAggregateID               = successor.HeaderAggregateID
	HeaderStreamVersion             = successor.HeaderStreamVersion
	HeaderEventName                 = successor.HeaderEventName
	HeaderEventSchemaVersion        = successor.HeaderEventSchemaVersion
	HeaderContentType               = successor.HeaderContentType
	HeaderRecordedAt                = successor.HeaderRecordedAt
	HeaderCorrelationID             = successor.HeaderCorrelationID
	HeaderCausationID               = successor.HeaderCausationID
	HeaderTenant                    = successor.HeaderTenant
	HeaderPartition                 = successor.HeaderPartition
	HeaderGlobalPosition            = successor.HeaderGlobalPosition
	HeaderApplicationMetadata       = successor.HeaderApplicationMetadata
	HeaderDeliveryMode              = successor.HeaderDeliveryMode
	HeaderDeadLetterSourceTopic     = successor.HeaderDeadLetterSourceTopic
	HeaderDeadLetterSourcePartition = successor.HeaderDeadLetterSourcePartition
	HeaderDeadLetterSourceOffset    = successor.HeaderDeadLetterSourceOffset
	HeaderDeadLetterSourceTime      = successor.HeaderDeadLetterSourceTime
	FailureRetry                    = successor.FailureRetry
	FailureHandled                  = successor.FailureHandled
)

var (
	ErrResolverRequired          = successor.ErrResolverRequired
	ErrResolverPanic             = successor.ErrResolverPanic
	ErrInvalidConfig             = successor.ErrInvalidConfig
	ErrTopicDenied               = successor.ErrTopicDenied
	ErrRecordInvalid             = successor.ErrRecordInvalid
	ErrRecordCorrupt             = successor.ErrRecordCorrupt
	ErrInvalidDeadLetterTopic    = successor.ErrInvalidDeadLetterTopic
	ErrInvalidDeadLetterConfig   = successor.ErrInvalidDeadLetterConfig
	ErrFailureCauseRequired      = successor.ErrFailureCauseRequired
	ErrDeadLetterLoop            = successor.ErrDeadLetterLoop
	ErrDeadLetterPublisherPanic  = successor.ErrDeadLetterPublisherPanic
	ErrDeadLetterPublishFailed   = successor.ErrDeadLetterPublishFailed
	ErrPublisherRequired         = successor.ErrPublisherRequired
	ErrCodecRequired             = successor.ErrCodecRequired
	ErrContextRequired           = successor.ErrContextRequired
	ErrReplayDenied              = successor.ErrReplayDenied
	ErrInvalidDispatcherOption   = successor.ErrInvalidDispatcherOption
	ErrPublisherPanic            = successor.ErrPublisherPanic
	ErrDispatchFailed            = successor.ErrDispatchFailed
	ErrConsumerRequired          = successor.ErrConsumerRequired
	ErrReplayHandlingDenied      = successor.ErrReplayHandlingDenied
	ErrInvalidHandlerOption      = successor.ErrInvalidHandlerOption
	ErrConsumerPanic             = successor.ErrConsumerPanic
	ErrRecordHandlingFailed      = successor.ErrRecordHandlingFailed
	ErrInvalidKafkaPosition      = successor.ErrInvalidKafkaPosition
	ErrFailurePolicyRequired     = successor.ErrFailurePolicyRequired
	ErrFailurePolicyPanic        = successor.ErrFailurePolicyPanic
	ErrInvalidFailureDisposition = successor.ErrInvalidFailureDisposition
)

type (
	DeadLetterError        = successor.DeadLetterError
	DeadLetterPolicy       = successor.DeadLetterPolicy
	DeadLetterPolicyConfig = successor.DeadLetterPolicyConfig
	DeliveryConsumer       = successor.DeliveryConsumer
	DeliveryConsumerFunc   = successor.DeliveryConsumerFunc
	DispatchError          = successor.DispatchError
	Dispatcher             = successor.Dispatcher
	DispatcherOption       = successor.DispatcherOption
	FailureDisposition     = successor.FailureDisposition
	FailurePolicy          = successor.FailurePolicy
	FailurePolicyFunc      = successor.FailurePolicyFunc
	HandlerError           = successor.HandlerError
	Publisher              = successor.Publisher
	RecordCodec            = successor.RecordCodec
	RecordCodecConfig      = successor.RecordCodecConfig
	RecordError            = successor.RecordError
	RecordHandler          = successor.RecordHandler
	RecordHandlerOption    = successor.RecordHandlerOption
	TopicResolver          = successor.TopicResolver
	TopicResolverFunc      = successor.TopicResolverFunc
)

func DefaultRecordLimits() broker.MessageLimits { return successor.DefaultRecordLimits() }

func NewDeadLetterPolicy(publisher Publisher, config DeadLetterPolicyConfig) (*DeadLetterPolicy, error) {
	return successor.NewDeadLetterPolicy(publisher, config)
}

func NewDispatcher(publisher Publisher, codec *RecordCodec, options ...DispatcherOption) (*Dispatcher, error) {
	return successor.NewDispatcher(publisher, codec, options...)
}

func AllowReplay() DispatcherOption { return successor.AllowReplay() }

func ContinueOnPublishError() DispatcherOption { return successor.ContinueOnPublishError() }

func AllowReplayHandling() RecordHandlerOption { return successor.AllowReplayHandling() }

func WithFailurePolicy(policy FailurePolicy) RecordHandlerOption {
	return successor.WithFailurePolicy(policy)
}

func NewRecordHandler(codec *RecordCodec, consumer DeliveryConsumer, options ...RecordHandlerOption) (*RecordHandler, error) {
	return successor.NewRecordHandler(codec, consumer, options...)
}

func NewRecordCodec(config RecordCodecConfig) (*RecordCodec, error) {
	return successor.NewRecordCodec(config)
}

func FixedTopic(topic string) TopicResolver { return successor.FixedTopic(topic) }
