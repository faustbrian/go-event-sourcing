package eventsourcing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestDecoderRejectsNilReceiverAndContextWithAssignedMessages(t *testing.T) {
	t.Parallel()

	message := Message{pending: invariantPendingMessage(), streamVersion: 1}
	var decoder *EventDecoder
	if events, err := decoder.DecodeContext(context.Background(), message); events != nil ||
		!errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("nil decoder DecodeContext() = %#v, %v", events, err)
	}
	decoder = &EventDecoder{}
	var nilContext context.Context
	if events, err := decoder.DecodeContext(nilContext, message); events != nil ||
		!errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("nil context DecodeContext() = %#v, %v", events, err)
	}
}

func TestCanonicalCharacterBoundariesRemainInclusive(t *testing.T) {
	t.Parallel()

	for _, value := range []string{"z", "az"} {
		if !validName(value, len(value)) {
			t.Fatalf("validName(%q) = false", value)
		}
	}
	for _, value := range []string{"a", "z", "A", "Z", "0", "9"} {
		if !validToken(value, 1) {
			t.Fatalf("validToken(%q) = false", value)
		}
	}
	if !validMetadataKey("es") || validMetadataKey("es.") {
		t.Fatal("metadata reserved-prefix boundary is incorrect")
	}
}

func TestJSONObjectDepthPropagatesAcrossObjectValues(t *testing.T) {
	t.Parallel()

	valid := strings.Repeat(`{"value":`, MaxJSONDepth) + "0" +
		strings.Repeat("}", MaxJSONDepth)
	if err := validateJSONValue(json.NewDecoder(strings.NewReader(valid)), 0); err != nil {
		t.Fatalf("validateJSONValue(maximum object depth) error = %v", err)
	}
	excessive := `{"value":` + valid + "}"
	if err := validateJSONValue(
		json.NewDecoder(strings.NewReader(excessive)),
		0,
	); !errors.Is(err, ErrMalformedEvent) {
		t.Fatalf("validateJSONValue(excessive object depth) error = %v", err)
	}
}

func TestJSONValidationRejectsUnterminatedContainer(t *testing.T) {
	t.Parallel()

	err := validateJSONValue(json.NewDecoder(strings.NewReader("[1}")), 0)
	if !errors.Is(err, ErrMalformedEvent) ||
		!strings.Contains(err.Error(), "unterminated JSON container") {
		t.Fatalf("validateJSONValue(unterminated array) error = %v", err)
	}
}

func TestDispatcherContinuesPastAFilteredConsumer(t *testing.T) {
	t.Parallel()

	delivery := Delivery{
		message: Message{pending: invariantPendingMessage(), streamVersion: 1},
		mode:    DeliveryLive,
	}
	first := Consumer{
		id:      "filtered",
		handler: func(context.Context, Delivery) error { t.Fatal("filtered handler called"); return nil },
		filters: []DeliveryFilter{func(Delivery) bool { return false }},
	}
	calls := 0
	second := Consumer{
		id: "selected",
		handler: func(context.Context, Delivery) error {
			calls++

			return nil
		},
	}
	dispatcher := &SyncDispatcher{consumers: []Consumer{first, second}}
	if err := dispatcher.Dispatch(context.Background(), []Delivery{delivery}); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("selected consumer calls = %d, want 1", calls)
	}
}

func TestTranslatorLimitAccountsForPriorStageOutput(t *testing.T) {
	t.Parallel()

	delivery := Delivery{
		message: Message{pending: invariantPendingMessage(), streamVersion: 1},
		mode:    DeliveryLive,
	}
	split := DeliveryTranslatorFunc(func(context.Context, Delivery) ([]Delivery, error) {
		return []Delivery{delivery, delivery}, nil
	})
	exact := DeliveryTranslatorFunc(func(context.Context, Delivery) ([]Delivery, error) {
		return []Delivery{delivery}, nil
	})
	chain, err := NewTranslatorChain(TranslatorChainConfig{
		Translators:   []DeliveryTranslator{split, exact},
		MaxDeliveries: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	translated, err := chain.Translate(context.Background(), delivery)
	if err != nil || len(translated) != 2 {
		t.Fatalf("Translate(exact bound) = %d, %v", len(translated), err)
	}

	overflow := DeliveryTranslatorFunc(func(context.Context, Delivery) ([]Delivery, error) {
		return []Delivery{delivery, delivery}, nil
	})
	chain, err = NewTranslatorChain(TranslatorChainConfig{
		Translators:   []DeliveryTranslator{split, overflow},
		MaxDeliveries: 3,
	})
	if err != nil {
		t.Fatal(err)
	}
	if translated, err = chain.Translate(context.Background(), delivery); translated != nil ||
		!errors.Is(err, ErrTranslationLimit) {
		t.Fatalf("Translate(accumulated overflow) = %#v, %v", translated, err)
	}
}

func TestUpcastLimitsAreInclusiveAndAccountForAccumulatedResults(t *testing.T) {
	t.Parallel()

	leaf := invariantUpcastEvent("event.leaf", 1)
	rule := UpcastRule{
		name:    EventName{value: "event.root"},
		version: 1,
		upcaster: func(UpcastEvent) ([]UpcastEvent, error) {
			return []UpcastEvent{leaf}, nil
		},
	}
	chain := &UpcasterChain{rules: map[upcastIdentity]UpcastRule{
		{name: "event.root", version: 1}: rule,
	}}
	root := invariantUpcastEvent("event.root", 1)
	path := map[upcastIdentity]struct{}{identityOf(root): {}}
	work := 4_095
	if output, err := chain.upcastOne(root, path, MaxUpcastSteps-1, &work); err != nil ||
		len(output) != 1 {
		t.Fatalf("upcastOne(exact work/step bound) = %d, %v", len(output), err)
	}
	work = 0
	if output, err := chain.upcastOne(root, path, MaxUpcastSteps, &work); output != nil ||
		!errors.Is(err, ErrUpcastLimit) {
		t.Fatalf("upcastOne(excessive steps) = %#v, %v", output, err)
	}
	work = 4_096
	if output, err := chain.upcastOne(root, path, 0, &work); output != nil ||
		!errors.Is(err, ErrUpcastLimit) {
		t.Fatalf("upcastOne(exhausted work) = %#v, %v", output, err)
	}

	exactSingle := singleExpansionUpcasterChain(MaxUpcastSegments)
	output, err := exactSingle.Upcast(invariantUpcastEvent("event.root", 1))
	if err != nil || len(output) != MaxUpcastSegments {
		t.Fatalf("Upcast(exact single expansion) = %d, %v", len(output), err)
	}
	excessiveSingle := singleExpansionUpcasterChain(MaxUpcastSegments + 1)
	if output, err = excessiveSingle.Upcast(
		invariantUpcastEvent("event.root", 1),
	); output != nil || !errors.Is(err, ErrUpcastLimit) {
		t.Fatalf("Upcast(excessive single expansion) = %#v, %v", output, err)
	}

	exact := accumulatedUpcasterChain(MaxUpcastSegments - 1)
	output, err = exact.Upcast(invariantUpcastEvent("event.root", 1))
	if err != nil || len(output) != MaxUpcastSegments {
		t.Fatalf("Upcast(exact accumulated bound) = %d, %v", len(output), err)
	}
	overflow := accumulatedUpcasterChain(MaxUpcastSegments)
	if output, err = overflow.Upcast(invariantUpcastEvent("event.root", 1)); output != nil ||
		!errors.Is(err, ErrUpcastLimit) {
		t.Fatalf("Upcast(accumulated overflow) = %#v, %v", output, err)
	}
}

func TestPlannedMessageValidationSeparatesIdentityAndVersion(t *testing.T) {
	t.Parallel()

	pending := invariantPendingMessage()
	plan := SavePlan{
		changes: ChangeSet{base: 4},
		prepared: []PendingMessage{
			pending,
			pending,
		},
	}
	messages := []Message{
		{pending: pending, streamVersion: 5},
		{pending: pending, streamVersion: 6},
	}
	if err := validatePlannedMessages(plan, messages); err != nil {
		t.Fatalf("validatePlannedMessages(valid) error = %v", err)
	}
	identityMismatch := cloneMessages(messages)
	identityMismatch[1].pending.id = MessageID{value: "different-message"}
	if err := validatePlannedMessages(plan, identityMismatch); !errors.Is(
		err,
		ErrPersistenceMismatch,
	) {
		t.Fatalf("validatePlannedMessages(identity mismatch) error = %v", err)
	}
	versionMismatch := cloneMessages(messages)
	versionMismatch[1].streamVersion = 5
	if err := validatePlannedMessages(plan, versionMismatch); !errors.Is(
		err,
		ErrPersistenceMismatch,
	) {
		t.Fatalf("validatePlannedMessages(version mismatch) error = %v", err)
	}
}

func TestDispatchCommittedValidatesReceiverAndContextIndependently(t *testing.T) {
	t.Parallel()

	owner := &repositorySaveOwner{}
	repository := &AggregateRepository[int, int]{saveOwner: owner}
	result := SaveResult{owner: owner, outcome: CommitCommitted}
	var nilContext context.Context
	if _, err := repository.DispatchCommitted(nilContext, result); !errors.Is(
		err,
		ErrInvalidArgument,
	) {
		t.Fatalf("DispatchCommitted(nil context) error = %v", err)
	}
	var nilRepository *AggregateRepository[int, int]
	if _, err := nilRepository.DispatchCommitted(context.Background(), result); !errors.Is(
		err,
		ErrInvalidArgument,
	) {
		t.Fatalf("nil DispatchCommitted() error = %v", err)
	}
}

func TestRepositoryPageCompletionChecksEachTerminationBoundary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		count     uint32
		batchSize uint32
		committed uint64
		want      bool
	}{
		{name: "exact full page continues", count: 2, batchSize: 2},
		{name: "short page stops", count: 1, batchSize: 2, want: true},
		{name: "maximum version stops", count: 2, batchSize: 2, committed: ^uint64(0), want: true},
	}
	for _, testCase := range tests {
		if got := repositoryPageComplete(
			testCase.count,
			testCase.batchSize,
			testCase.committed,
		); got != testCase.want {
			t.Fatalf("%s: repositoryPageComplete() = %t, want %t", testCase.name, got, testCase.want)
		}
	}
}

func invariantUpcastEvent(name string, version SchemaVersion) UpcastEvent {
	event := invariantEncodedEvent()
	event.name = EventName{value: name}
	event.version = version

	return UpcastEvent{event: event}
}

func accumulatedUpcasterChain(firstExpansion int) *UpcasterChain {
	root := invariantUpcastEvent("event.root", 1)
	branch := invariantUpcastEvent("event.branch", 1)
	tail := invariantUpcastEvent("event.tail", 1)
	leaves := make([]UpcastEvent, firstExpansion)
	for index := range leaves {
		leaves[index] = invariantUpcastEvent(fmt.Sprintf("event.leaf-%d", index), 1)
	}
	rules := map[upcastIdentity]UpcastRule{
		identityOf(root): {
			name:    root.event.name,
			version: root.event.version,
			upcaster: func(UpcastEvent) ([]UpcastEvent, error) {
				return []UpcastEvent{branch, tail}, nil
			},
		},
		identityOf(branch): {
			name:    branch.event.name,
			version: branch.event.version,
			upcaster: func(UpcastEvent) ([]UpcastEvent, error) {
				return leaves, nil
			},
		},
	}

	return &UpcasterChain{rules: rules}
}

func singleExpansionUpcasterChain(expansion int) *UpcasterChain {
	root := invariantUpcastEvent("event.root", 1)
	leaves := make([]UpcastEvent, expansion)
	for index := range leaves {
		leaves[index] = invariantUpcastEvent(fmt.Sprintf("event.leaf-%d", index), 1)
	}

	return &UpcasterChain{rules: map[upcastIdentity]UpcastRule{
		identityOf(root): {
			name:    root.event.name,
			version: root.event.version,
			upcaster: func(UpcastEvent) ([]UpcastEvent, error) {
				return leaves, nil
			},
		},
	}}
}
