# Adapter path migration

## Core-v2 adoption

The upcoming canonical Kafka-v2 and OpenTelemetry-v2 adapters adopt public
core-v2. After their independent releases are published, require each
`/adapters/<name>/v2` module and migrate caller core imports/types to
`github.com/faustbrian/go-event-sourcing/v2` together. See the
[Kafka](../adapters/kafka/docs/reference.md#v2-migration) and
[OpenTelemetry](../adapters/otel/docs/reference.md#v2-migration) guides.

The following selector-preserving migration describes the published v1
successors for applications retaining core-v1. Deprecated gokafka-v1 and
gotelemetry-v1 remain unchanged; core-v2 availability does not retire them.

## Published v1 path migration

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

The deprecated modules retain their own released v1 implementations, concrete
types, sentinel errors, reflection paths, and `%T` output for the frozen support
interval. The successors preserve the API shape, Kafka wire behavior, telemetry
signals, and runtime contracts, but migration changes package-qualified type,
reflection, and sentinel error identities. Update type assertions and
`errors.Is` comparisons to use the selected import path consistently.

## Release order

1. Release `adapters/kafka/v1.0.0` and `adapters/otel/v1.0.0`.
2. Verify both modules through clean public module resolution.
3. Verify the deprecated modules against their immutable v1 API baselines and
   the public successor behavior.
4. Release new patch versions of `adapters/gokafka` and
   `adapters/gotelemetry`.

This order makes both migration targets publicly resolvable before their
predecessors advertise them. The successor implementations do not import either
deprecated path.

## Compatibility interval

Each released v1 path remains supported for the longer of 180 days and two
published stable minor releases after its successor is publicly consumable.
Removal requires an authorized next-major release and a fresh consumer audit.
The interval does not permit silent behavior, wire, or telemetry-scope drift.

## Fixture supplier adoption

The Testcontainers core v0.44.0 and Moby client v0.6.1 updates affect the
dependency closure of gokafka, gotelemetry, kafka, otel, outbox, queue, and
postgres. Kafka fixture modules use v0.44.0; PostgreSQL fixture modules
remain at v0.43.0 over the newer core. These container clients are used by
fixtures, not by the owned production adapter implementations. Owned APIs,
wire formats, transaction boundaries, and telemetry-provider ownership are
unchanged.

The selected Moby API module is v1.56.1. Client v0.6.1 automatically
negotiates Docker API versions 1.40 through 1.56; fixed API options or
environment overrides disable negotiation. The API module patch version
is not a Docker wire-version identifier.

Go minimum-version selection also selects `otelhttp` v0.69.0 in each module.
Applications importing that upstream package directly must account for its
removed APIs even when they do not run these fixtures:

- Read the `otelhttp.Version` constant instead of calling `Version()`.
- Replace `DefaultClient`, `Get`, `Head`, `Post`, and `PostForm` with an
  application-owned `http.Client` using `otelhttp.NewTransport`.
- Replace `WithPublicEndpoint` with `WithPublicEndpointFn`, preserving the
  application's public-endpoint policy.
- Remove `WithRouteTag`; routes are added to spans automatically. For
  application-owned metric attributes, use the current `Labeler` API rather
  than adding a dependency on the now-deprecated `WithMetricAttributesFn`.

See the [upstream changelog](https://github.com/open-telemetry/opentelemetry-go-contrib/blob/v1.44.0/CHANGELOG.md)
for the supplier's migration and HTTP semantic-convention changes. The
selected OpenTelemetry API and SDK remain v1.45.0; this adoption does not
upgrade them to v1.47.0 or change the event-sourcing module paths.
