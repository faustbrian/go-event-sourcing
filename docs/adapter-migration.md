# Adapter path migration

Target-oriented adapter paths replace the two released names that redundantly
include `go`:

| Released v1 path | Preferred successor | Migration |
| --- | --- | --- |
| `adapters/gokafka` | `adapters/kafka` | Change the import path and rename the qualifier to `kafka`, or alias the successor import as `gokafka`. |
| `adapters/gotelemetry` | `adapters/otel` | Change the import path and rename the qualifier to `otel`, or alias the successor import as `gotelemetry`. |

The successors preserve the complete public APIs and runtime contracts. Kafka
wire version, headers, record ownership, acknowledgements, replay policy,
errors, and failure dispositions do not change. OpenTelemetry signal names,
attributes, semantic convention, instrumentation scope, propagation, privacy
bounds, errors, and provider ownership do not change.

For a selector-preserving migration, use explicit compatibility aliases:

```go
import (
	gokafka "github.com/faustbrian/go-event-sourcing/adapters/kafka"
	gotelemetry "github.com/faustbrian/go-event-sourcing/adapters/otel"
)
```

The successor modules compile external-package migration fixtures with these
aliases. Applications may instead adopt the target-oriented `kafka` and `otel`
qualifiers and update selectors mechanically.

The deprecated modules implement their facades with Go type aliases. Source
selectors, methods, errors, Kafka wire behavior, and telemetry signals remain
compatible, but reflection and `%T` report the successor package path. Do not
use concrete Go package identity as a behavior or persistence key.

## Release order

1. Release `adapters/kafka/v1.0.0` and `adapters/otel/v1.0.0`.
2. Verify both modules through clean public module resolution.
3. Add those released versions to the legacy facade module requirements.
4. Release new patch versions of `adapters/gokafka` and
   `adapters/gotelemetry`.

This order prevents either released facade from naming an unavailable module.
The successor implementations do not import either deprecated path.

## Compatibility interval

Each released v1 path remains supported for the longer of 180 days and two
published stable minor releases after its successor is publicly consumable.
Removal requires an authorized next-major release and a fresh consumer audit.
The interval does not permit silent behavior, wire, or telemetry-scope drift.
