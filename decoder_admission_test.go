package eventsourcing_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	eventsourcing "github.com/faustbrian/go-event-sourcing/v2"
)

type decoderSizedUpcaster int

func (size decoderSizedUpcaster) Upcast(event eventsourcing.UpcastEvent) ([]eventsourcing.UpcastEvent, error) {
	output := make([]eventsourcing.UpcastEvent, int(size))
	for index := range output {
		output[index] = event
	}
	return output, nil
}

type decoderContextSizedUpcaster struct{ decoderSizedUpcaster }

func (upcaster decoderContextSizedUpcaster) UpcastContext(_ context.Context, event eventsourcing.UpcastEvent) ([]eventsourcing.UpcastEvent, error) {
	return upcaster.decoderSizedUpcaster.Upcast(event)
}

func TestEventDecoderBoundsCustomUpcasterSegments(t *testing.T) {
	stream, err := eventsourcing.NewStreamID("account", "decoder-limit")
	if err != nil {
		t.Fatal(err)
	}
	message := persistedRepositoryMessage(t, stream,
		mustEncodedEvent(t, "account.opened", 1, []byte(`{"owner":"Ada"}`)))
	for _, contextual := range []bool{false, true} {
		for _, count := range []int{0, eventsourcing.MaxUpcastSegments, eventsourcing.MaxUpcastSegments + 1} {
			t.Run(fmt.Sprintf("context=%t/count=%d", contextual, count), func(t *testing.T) {
				var upcaster eventsourcing.Upcaster = decoderSizedUpcaster(count)
				if contextual {
					upcaster = decoderContextSizedUpcaster{decoderSizedUpcaster(count)}
				}
				codec := repositoryCodec(t)
				decodeCalls := 0
				decoder, err := eventsourcing.NewEventDecoder(decoderCodec{decode: func(event eventsourcing.EncodedEvent) (eventsourcing.DecodedEvent, error) {
					decodeCalls++
					return codec.Decode(event)
				}}, upcaster)
				if err != nil {
					t.Fatal(err)
				}
				logical, err := decoder.DecodeContext(context.Background(), message)
				if count > eventsourcing.MaxUpcastSegments {
					if !errors.Is(err, eventsourcing.ErrUpcastLimit) || logical != nil || decodeCalls != 0 {
						t.Fatalf("excessive output: events=%d error=%v decode calls=%d", len(logical), err, decodeCalls)
					}
					return
				}
				if err != nil || len(logical) != count || decodeCalls != count {
					t.Fatalf("bounded output: events=%d error=%v decode calls=%d", len(logical), err, decodeCalls)
				}
				for index, event := range logical {
					if event.SegmentIndex() != uint32(index) || event.SegmentCount() != uint32(count) || event.IsZero() {
						t.Fatalf("incorrect segment coordinates at %d", index)
					}
				}
			})
		}
	}
}
