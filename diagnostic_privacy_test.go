package eventsourcing_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	eventsourcing "github.com/faustbrian/go-event-sourcing"
)

func TestJSONCodecDiagnosticsDoNotDiscloseInput(t *testing.T) {
	const privateName = "private.fixture.event"
	const privateType = "application/private-fixture"
	codec, err := eventsourcing.NewJSONCodec(eventsourcing.JSONEvent[struct{}](privateName, 1))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name        string
		eventName   string
		version     eventsourcing.SchemaVersion
		contentType string
		payload     string
		category    error
	}{
		{"content type", privateName, 1, privateType, `{}`, eventsourcing.ErrUnsupportedContentType},
		{"unknown identity", "private.fixture.unknown", 1, eventsourcing.JSONContentType, `{}`, eventsourcing.ErrUnknownEvent},
		{"incompatible version", privateName, 2, eventsourcing.JSONContentType, `{}`, eventsourcing.ErrIncompatibleVersion},
		{"malformed payload", privateName, 1, eventsourcing.JSONContentType, `{"private.fixture.payload":`, eventsourcing.ErrMalformedEvent},
	} {
		t.Run(test.name, func(t *testing.T) {
			event, err := eventsourcing.NewEncodedEvent(eventsourcing.EncodedEventInput{
				Name: test.eventName, Version: test.version, ContentType: test.contentType, Payload: []byte(test.payload),
			})
			if err != nil {
				t.Fatal(err)
			}
			_, err = codec.Decode(event)
			assertPrivateDiagnostic(t, err, test.category, test.eventName, privateType, "private.fixture.payload")
		})
	}
	for _, test := range []struct {
		name     string
		version  eventsourcing.SchemaVersion
		category error
	}{
		{"private.fixture.unknown", 1, eventsourcing.ErrUnknownEvent},
		{privateName, 2, eventsourcing.ErrIncompatibleVersion},
		{privateName, 1, eventsourcing.ErrEventTypeMismatch},
	} {
		t.Run("encode "+test.name+" "+test.category.Error(), func(t *testing.T) {
			event, err := eventsourcing.NewDecodedEvent(eventsourcing.DecodedEventInput{
				Name: test.name, Version: test.version, Value: "private.fixture.value",
			})
			if err != nil {
				t.Fatal(err)
			}
			_, err = codec.Encode(event)
			assertPrivateDiagnostic(t, err, test.category, test.name, "private.fixture.value")
		})
	}
}

func TestRegistrationDiagnosticsDoNotDiscloseIdentities(t *testing.T) {
	const privateName = "private.fixture.event"
	const privateAlias = "private.fixture.alias"
	registration := eventsourcing.JSONEvent[struct{}](privateName, 1)
	for _, test := range []struct {
		name     string
		options  []eventsourcing.JSONCodecOption
		category error
	}{
		{"duplicate event", []eventsourcing.JSONCodecOption{registration, registration}, eventsourcing.ErrDuplicateRegistration},
		{"duplicate alias", []eventsourcing.JSONCodecOption{registration, eventsourcing.JSONAlias(privateAlias, 1, privateName, 1), eventsourcing.JSONAlias(privateAlias, 1, privateName, 1)}, eventsourcing.ErrDuplicateRegistration},
		{"missing alias target", []eventsourcing.JSONCodecOption{registration, eventsourcing.JSONAlias(privateAlias, 1, "private.fixture.missing", 1)}, eventsourcing.ErrUnknownEvent},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := eventsourcing.NewJSONCodec(test.options...)
			assertPrivateDiagnostic(t, err, test.category, privateName, privateAlias, "private.fixture.missing")
		})
	}
	rule, err := eventsourcing.NewUpcastRule(privateName, 1, func(event eventsourcing.UpcastEvent) ([]eventsourcing.UpcastEvent, error) {
		return []eventsourcing.UpcastEvent{event}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = eventsourcing.NewUpcasterChain(rule, rule)
	assertPrivateDiagnostic(t, err, eventsourcing.ErrDuplicateRegistration, privateName)
	consumer, err := eventsourcing.NewConsumer("private-fixture-consumer", func(context.Context, eventsourcing.Delivery) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	_, err = eventsourcing.NewSyncDispatcher(consumer, consumer)
	assertPrivateDiagnostic(t, err, eventsourcing.ErrDuplicateConsumer, consumer.ID())
}

func assertPrivateDiagnostic(t *testing.T, err, category error, markers ...string) {
	t.Helper()
	if !errors.Is(err, category) {
		t.Fatalf("diagnostic category = %v, want %v", err, category)
	}
	for _, marker := range markers {
		if strings.Contains(err.Error(), marker) {
			t.Errorf("diagnostic disclosed fixture input %q", marker)
		}
	}
}
