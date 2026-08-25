package eventsourcing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestZeroSentinelsRequireEveryFieldToBeUnassigned(t *testing.T) {
	t.Parallel()

	if !(StreamID{}).IsZero() ||
		(StreamID{aggregateType: "account"}).IsZero() ||
		(StreamID{aggregateID: "42"}).IsZero() {
		t.Fatal("StreamID zero sentinel accepted a partially assigned value")
	}

	encodedParts := []EncodedEvent{
		{name: EventName{value: "account.opened"}},
		{version: 1},
		{contentType: JSONContentType},
		{payload: []byte(`{}`)},
	}
	if !(EncodedEvent{}).IsZero() {
		t.Fatal("zero EncodedEvent is assigned")
	}
	for index, event := range encodedParts {
		if event.IsZero() {
			t.Fatalf("partially assigned EncodedEvent %d is zero", index)
		}
	}

	decodedParts := []DecodedEvent{
		{name: EventName{value: "account.opened"}},
		{version: 1},
		{value: struct{}{}},
	}
	if !(DecodedEvent{}).IsZero() {
		t.Fatal("zero DecodedEvent is assigned")
	}
	for index, event := range decodedParts {
		if event.IsZero() {
			t.Fatalf("partially assigned DecodedEvent %d is zero", index)
		}
	}

	now := time.Date(2026, time.August, 25, 12, 0, 0, 0, time.UTC)
	snapshotParts := []Snapshot{
		{stream: StreamID{aggregateType: "account", aggregateID: "42"}},
		{aggregateVersion: 1},
		{schemaVersion: 1},
		{state: []byte(`{}`)},
		{metadata: map[string]string{"source": "test"}},
		{createdAt: now},
	}
	if !(Snapshot{}).IsZero() {
		t.Fatal("zero Snapshot is assigned")
	}
	for index, snapshot := range snapshotParts {
		if snapshot.IsZero() {
			t.Fatalf("partially assigned Snapshot %d is zero", index)
		}
	}
}

func TestRequiredEnvelopeComponentsFailClosedIndependently(t *testing.T) {
	t.Parallel()

	pending := invariantPendingMessage()
	message := Message{pending: pending, streamVersion: 1, globalPosition: 1}
	decoded := invariantDecodedEvent()
	logical := LogicalEvent{
		source:       message,
		event:        decoded,
		segmentCount: 1,
	}
	if logical.IsZero() {
		t.Fatal("fully assigned LogicalEvent is zero")
	}
	logicalCases := []LogicalEvent{
		{event: decoded, segmentCount: 1},
		{source: message, segmentCount: 1},
		{source: message, event: decoded},
	}
	for index, event := range logicalCases {
		if !event.IsZero() {
			t.Fatalf("incomplete LogicalEvent %d is assigned", index)
		}
	}

	if !validPendingMessage(pending) {
		t.Fatal("fully assigned PendingMessage is invalid")
	}
	pendingCases := []PendingMessage{
		pending,
		pending,
		pending,
		pending,
	}
	pendingCases[0].id = MessageID{}
	pendingCases[1].stream = StreamID{}
	pendingCases[2].event = EncodedEvent{}
	pendingCases[3].recordedAt = time.Time{}
	for index, candidate := range pendingCases {
		if validPendingMessage(candidate) {
			t.Fatalf("incomplete PendingMessage %d is valid", index)
		}
		if _, err := candidate.WithMetadata(nil); !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("WithMetadata(incomplete %d) error = %v", index, err)
		}
		if _, err := NewMessage(MessageInput{
			Pending:       candidate,
			StreamVersion: 1,
		}); !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("NewMessage(incomplete %d) error = %v", index, err)
		}
	}

	validDelivery := Delivery{message: message, mode: DeliveryLive}
	if validDelivery.IsZero() {
		t.Fatal("fully assigned Delivery is zero")
	}
	for index, delivery := range []Delivery{
		{mode: DeliveryLive},
		{message: message},
	} {
		if !delivery.IsZero() {
			t.Fatalf("incomplete Delivery %d is assigned", index)
		}
	}
}

func TestValueEqualityObservesEveryIndependentField(t *testing.T) {
	t.Parallel()

	base := Message{
		pending:        invariantPendingMessage(),
		streamVersion:  3,
		globalPosition: 7,
	}
	if !base.Equal(base) {
		t.Fatal("Message is not equal to itself")
	}
	variants := []Message{base, base, base}
	variants[0].pending.id = MessageID{value: "other-message"}
	variants[1].streamVersion++
	variants[2].globalPosition++
	for index, variant := range variants {
		if base.Equal(variant) {
			t.Fatalf("Message equality ignored field variant %d", index)
		}
	}

	upcast := UpcastEvent{
		event:    invariantEncodedEvent(),
		metadata: map[string]string{"source": "test"},
	}
	if !upcastEventsEqual([]UpcastEvent{upcast}, []UpcastEvent{upcast}) {
		t.Fatal("equivalent upcast events differ")
	}
	upcastVariants := []UpcastEvent{upcast, upcast, upcast, upcast, upcast}
	upcastVariants[0].event.name = EventName{value: "account.closed"}
	upcastVariants[1].event.version++
	upcastVariants[2].event.contentType = "application/octet-stream"
	upcastVariants[3].event.payload = []byte(`{"different":true}`)
	upcastVariants[4].metadata = map[string]string{"source": "other"}
	for index, variant := range upcastVariants {
		if upcastEventsEqual([]UpcastEvent{upcast}, []UpcastEvent{variant}) {
			t.Fatalf("upcast equality ignored field variant %d", index)
		}
	}
}

func TestInternalRangeValidityRejectsEachMalformedField(t *testing.T) {
	t.Parallel()

	stream := ReadStreamOptions{fromVersion: 2, toVersion: 2, limit: MaxReadMessages}
	if !stream.Valid() {
		t.Fatal("inclusive maximum ReadStreamOptions is invalid")
	}
	streamCases := []ReadStreamOptions{stream, stream, stream, stream}
	streamCases[0].fromVersion = 0
	streamCases[1].toVersion = 1
	streamCases[2].limit = 0
	streamCases[3].limit = MaxReadMessages + 1
	for index, options := range streamCases {
		if options.Valid() {
			t.Fatalf("malformed ReadStreamOptions %d is valid", index)
		}
	}

	global := ReadGlobalOptions{
		fromPosition: 2,
		toPosition:   2,
		limit:        MaxReadMessages,
	}
	if !global.Valid() {
		t.Fatal("inclusive maximum ReadGlobalOptions is invalid")
	}
	globalCases := []ReadGlobalOptions{global, global, global, global}
	globalCases[0].fromPosition = 0
	globalCases[1].toPosition = 1
	globalCases[2].limit = 0
	globalCases[3].limit = MaxReadMessages + 1
	for index, options := range globalCases {
		if options.Valid() {
			t.Fatalf("malformed ReadGlobalOptions %d is valid", index)
		}
	}
}

func TestScalarValidatorsExerciseEveryCanonicalBoundary(t *testing.T) {
	t.Parallel()

	for _, value := range []string{"a", "a.a", "a0", "a-", "a_"} {
		if !validName(value, len(value)) {
			t.Fatalf("validName(%q) = false", value)
		}
	}
	for _, value := range []string{"", "A", "0a", ".a", "a.", "a..a", "a/A"} {
		if validName(value, MaxEventNameBytes) {
			t.Fatalf("validName(%q) = true", value)
		}
	}

	for _, value := range []string{"a", "A", "0", ".", "_", ":", "-"} {
		if !validToken(value, 1) {
			t.Fatalf("validToken(%q) = false", value)
		}
	}
	for _, value := range []string{"", "/", " "} {
		if validToken(value, MaxMessageIDBytes) {
			t.Fatalf("validToken(%q) = true", value)
		}
	}

	if !validText(strings.Repeat("x", 8), 8, false) ||
		validText(strings.Repeat("x", 9), 8, false) ||
		validText(" ", 8, false) ||
		!validText("", 8, true) ||
		validText(string([]byte{0xff}), 8, true) ||
		validText("line\nfeed", 32, true) {
		t.Fatal("validText did not enforce exact text boundaries")
	}

	maximumContentType := "application/" +
		strings.Repeat("a", MaxContentTypeBytes-len("application/"))
	if !validContentType(maximumContentType) ||
		validContentType(maximumContentType+"a") ||
		validContentType("application/json/extra") {
		t.Fatal("validContentType did not enforce exact media-type boundaries")
	}
}

func TestPayloadMetadataAndSnapshotBoundsAreInclusive(t *testing.T) {
	t.Parallel()

	stream := StreamID{aggregateType: "account", aggregateID: "42"}
	for _, size := range []int{1, MaxPayloadBytes} {
		if _, err := NewEncodedEvent(EncodedEventInput{
			Name:        "account.opened",
			Version:     1,
			ContentType: JSONContentType,
			Payload:     make([]byte, size),
		}); err != nil {
			t.Fatalf("NewEncodedEvent(payload=%d) error = %v", size, err)
		}
	}
	for _, size := range []int{0, MaxPayloadBytes + 1} {
		if _, err := NewEncodedEvent(EncodedEventInput{
			Name:        "account.opened",
			Version:     1,
			ContentType: JSONContentType,
			Payload:     make([]byte, size),
		}); !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("NewEncodedEvent(payload=%d) error = %v", size, err)
		}
	}

	metadata := exactSizeMetadata()
	if _, err := copyMetadata(metadata); err != nil {
		t.Fatalf("copyMetadata(exact maximum) error = %v", err)
	}
	metadata["key-00"] += "x"
	if _, err := copyMetadata(metadata); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("copyMetadata(over maximum) error = %v", err)
	}

	now := time.Date(2026, time.August, 25, 12, 0, 0, 0, time.UTC)
	for _, size := range []int{1, MaxSnapshotStateBytes} {
		if _, err := NewSnapshot(SnapshotInput{
			Stream:           stream,
			AggregateVersion: 1,
			SchemaVersion:    1,
			State:            make([]byte, size),
			CreatedAt:        now,
		}); err != nil {
			t.Fatalf("NewSnapshot(state=%d) error = %v", size, err)
		}
	}
	for _, size := range []int{0, MaxSnapshotStateBytes + 1} {
		if _, err := NewSnapshot(SnapshotInput{
			Stream:           stream,
			AggregateVersion: 1,
			SchemaVersion:    1,
			State:            make([]byte, size),
			CreatedAt:        now,
		}); !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("NewSnapshot(state=%d) error = %v", size, err)
		}
	}
}

func TestJSONStructuralLimitsFailAtTheFirstExcessValue(t *testing.T) {
	t.Parallel()

	validDepth := strings.Repeat("[", MaxJSONDepth) + "0" +
		strings.Repeat("]", MaxJSONDepth)
	if err := validateJSONValue(json.NewDecoder(strings.NewReader(validDepth)), 0); err != nil {
		t.Fatalf("validateJSONValue(maximum depth) error = %v", err)
	}
	excessiveDepth := "[" + validDepth + "]"
	if err := validateJSONValue(
		json.NewDecoder(strings.NewReader(excessiveDepth)),
		0,
	); !errors.Is(err, ErrMalformedEvent) {
		t.Fatalf("validateJSONValue(excessive depth) error = %v", err)
	}

	validArray := "[" + strings.Repeat("0,", MaxJSONContainerEntries-1) + "0]"
	if err := validateJSONValue(json.NewDecoder(strings.NewReader(validArray)), 0); err != nil {
		t.Fatalf("validateJSONValue(maximum array) error = %v", err)
	}
	excessiveArray := "[" + strings.Repeat("0,", MaxJSONContainerEntries) + "0]"
	if err := validateJSONValue(
		json.NewDecoder(strings.NewReader(excessiveArray)),
		0,
	); !errors.Is(err, ErrMalformedEvent) {
		t.Fatalf("validateJSONValue(excessive array) error = %v", err)
	}

	validObject := jsonObjectWithEntries(MaxJSONContainerEntries)
	if err := validateJSONValue(json.NewDecoder(strings.NewReader(validObject)), 0); err != nil {
		t.Fatalf("validateJSONValue(maximum object) error = %v", err)
	}
	excessiveObject := jsonObjectWithEntries(MaxJSONContainerEntries + 1)
	if err := validateJSONValue(
		json.NewDecoder(strings.NewReader(excessiveObject)),
		0,
	); !errors.Is(err, ErrMalformedEvent) {
		t.Fatalf("validateJSONValue(excessive object) error = %v", err)
	}
}

func TestConstructedRulesRejectEachMissingComponent(t *testing.T) {
	t.Parallel()

	rule := UpcastRule{
		name:     EventName{value: "account.opened"},
		version:  1,
		upcaster: func(event UpcastEvent) ([]UpcastEvent, error) { return []UpcastEvent{event}, nil },
	}
	for index, candidate := range []UpcastRule{
		{version: rule.version, upcaster: rule.upcaster},
		{name: rule.name, upcaster: rule.upcaster},
		{name: rule.name, version: rule.version},
	} {
		if _, err := NewUpcasterChain(candidate); !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("NewUpcasterChain(incomplete %d) error = %v", index, err)
		}
	}

	consumer := Consumer{id: "projection", handler: func(context.Context, Delivery) error {
		return nil
	}}
	for index, candidate := range []Consumer{
		{handler: consumer.handler},
		{id: consumer.id},
	} {
		builder := syncDispatcherBuilder{identities: map[string]struct{}{}}
		if err := candidate.configureSyncDispatcher(&builder); !errors.Is(
			err,
			ErrInvalidArgument,
		) {
			t.Fatalf("configureSyncDispatcher(incomplete %d) error = %v", index, err)
		}
	}
}

func invariantEncodedEvent() EncodedEvent {
	return EncodedEvent{
		name:        EventName{value: "account.opened"},
		version:     1,
		contentType: JSONContentType,
		payload:     []byte(`{}`),
	}
}

func invariantDecodedEvent() DecodedEvent {
	return DecodedEvent{
		name:    EventName{value: "account.opened"},
		version: 1,
		value:   struct{}{},
	}
}

func invariantPendingMessage() PendingMessage {
	return PendingMessage{
		id:         MessageID{value: "message-1"},
		stream:     StreamID{aggregateType: "account", aggregateID: "42"},
		event:      invariantEncodedEvent(),
		metadata:   map[string]string{"source": "test"},
		recordedAt: time.Date(2026, time.August, 25, 12, 0, 0, 0, time.UTC),
	}
}

func exactSizeMetadata() map[string]string {
	metadata := make(map[string]string, MaxMetadataEntries)
	total := 0
	for index := 0; index < MaxMetadataEntries; index++ {
		key := fmt.Sprintf("key-%02d", index)
		valueLength := 1_000
		if index == MaxMetadataEntries-1 {
			valueLength = MaxMetadataBytes - total - len(key)
		}
		metadata[key] = strings.Repeat("x", valueLength)
		total += len(key) + valueLength
	}

	return metadata
}

func jsonObjectWithEntries(entries int) string {
	var builder strings.Builder
	builder.WriteByte('{')
	for index := 0; index < entries; index++ {
		if index != 0 {
			builder.WriteByte(',')
		}
		fmt.Fprintf(&builder, "%q:0", fmt.Sprintf("key-%d", index))
	}
	builder.WriteByte('}')

	return builder.String()
}
