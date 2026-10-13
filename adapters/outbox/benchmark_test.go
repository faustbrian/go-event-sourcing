package eventoutbox

// This benchmark measures encode/decode round trips for a prebuilt complete
// message and fixed topic. Codec/fixture setup precedes ResetTimer; SQL, outer
// commit and broker publication are not measured.

import (
	"testing"

	"github.com/faustbrian/go-transactional-outbox/v2"
)

func BenchmarkEnvelopeCodecRoundTrip(b *testing.B) {
	codec := mustCodec(b, FixedTopic("events"), outbox.DefaultLimits())
	message := internalMessage(b, true)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		envelope, err := codec.Encode(message)
		if err != nil {
			b.Fatal(err)
		}
		if _, err := codec.Decode(envelope); err != nil {
			b.Fatal(err)
		}
	}
}
