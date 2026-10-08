# Release notes and compatibility

The changelog in each independently releasable module is the authoritative
release-note source for that module:

- [core event sourcing](../CHANGELOG.md);
- [PostgreSQL](../postgres/CHANGELOG.md);
- [Kafka](../adapters/kafka/CHANGELOG.md);
- [queue](../adapters/queue/CHANGELOG.md);
- [outbox](../adapters/outbox/CHANGELOG.md); and
- [OpenTelemetry](../adapters/otel/CHANGELOG.md);
- [deprecated Kafka v1 compatibility module](../adapters/gokafka/CHANGELOG.md); and
- [deprecated OpenTelemetry v1 compatibility module](../adapters/gotelemetry/CHANGELOG.md).

The published v1 surfaces follow semantic versioning. Unreleased changelog
entries describe pending changes and do not alter the latest tagged contract.
A release note must not imply that a local check, adapter test, or benchmark is
deployed or production-verified.

## Adopting core v2

For published v2.0.0, require
`github.com/faustbrian/go-event-sourcing/v2@v2.0.0` and insert `/v2` immediately
after `go-event-sourcing` in imports of the core, `eventtest`, `memory`,
`processmanager`, `projection`, and `snapshot`. Go 1.27.0 remains the minimum.
Source stays at the repository root, without a `v2` source directory.

The new module identity changes Go named types, interface method signatures,
sentinel identities, and reflection paths. Do not mix v1 and v2 messages,
stores, codecs, dispatchers, or projection interfaces. Stored event names and
schema versions remain explicit data; this import migration does not change
stored envelopes, wire formats, PostgreSQL schemas, or application history.

Built-in codec, registration, upcaster, and dispatcher errors no longer print
event identities or content types. Use `errors.Is` for stable categories and
explicit trusted cause inspection where documented, rather than parsing error
text. Application-owned callback errors remain a separate caller-controlled
boundary; this is not blanket callback redaction.

Independent nested modules retain their existing source locations and releases.
Published PostgreSQL-v2 adopts public core-v2. The published outbox-v2 major
adopts both public producers together; the other nested modules retain
their independent releases. Published Kafka-v2, queue-v2, and OpenTelemetry-v2
adopt public core-v2 directly; deprecated gokafka-v1 and gotelemetry-v1 retain
core-v1. The private competitor module adopts core-v2 without a public tag.
When a maintained module adopts core-v2 types, its
own major suffix follows its directory: for example,
`github.com/faustbrian/go-event-sourcing/postgres/v2`, not
`github.com/faustbrian/go-event-sourcing/v2/postgres`. Publish core-v2 first,
then PostgreSQL-v2 before an outbox-v2 adapter that depends on both. Queue,
Kafka, and OpenTelemetry adoption may follow core independently.

For those three published adapters, require their own `/adapters/<name>/v2`
identity and migrate caller core nominal types
in the same application change. Their other dependencies, wire formats,
algorithms, and telemetry instrumentation scope remain unchanged. See the
[Kafka](../adapters/kafka/docs/reference.md#v2-migration),
[queue](../adapters/queue/docs/reference.md#v2-migration), and
[OpenTelemetry](../adapters/otel/docs/reference.md#v2-migration) guides.

PostgreSQL-v2 changes public message, stream, version, snapshot, projection,
and prepared-plan identities without changing SQL or stored data. Its existing
`adapters/outbox` v1 consumer remains on PostgreSQL-v1; it needs a separate
major migration, not an automatic dependency upgrade. The published
`adapters/outbox/v2` changes the adapter's nominal API while preserving
Transactional Outbox v1 and all staging behavior. See its
[migration guide](../adapters/outbox/docs/reference.md#v2-migration) and the
[PostgreSQL migration guide](../postgres/docs/reference.md#v2-migration).

The deprecated `adapters/gokafka` and `adapters/gotelemetry` v1 implementations
remain unchanged during their existing support interval. Root-v2 availability
does not retire them or change their published contracts. Existing v1
applications can continue using their selected public releases.

`api/v1.0.0.txt` preserves the exact released root API from commit
`9f7bffbb757620e0e0796d9fb564bf66e3cabf7a`. `api/baseline.txt` is the current
root-major API snapshot; cross-major compatibility is an adoption decision,
not an assertion that the two module identities are interchangeable.

## Patch upgrade requirements

For custom-upcaster segment admission, select core
`github.com/faustbrian/go-event-sourcing/v2@v2.0.1` or newer explicitly.
This also applies when a maintained v2 adapter supplies an upcaster, including
the OpenTelemetry wrapper; an adapter's existing core-v2.0.0 minimum does not
itself supply the fix. Both legacy and context-aware upcasters keep the same
inclusive segment limit and empty-output behavior.

For absent nullable projection checkpoints, select
`github.com/faustbrian/go-event-sourcing/postgres/v2@v2.0.2` or newer.
The fix preserves valid status when an absent checkpoint carries an inactive
integer; it does not change SQL, stored data, or transaction ownership.
These patch targets retain the existing v2 type identities and Go 1.27 minimum;
deprecated core-v1 adapters remain separate.

## Semantic-versioning surfaces

Treat the following as compatibility-sensitive even when their Go declarations
do not change:

- event names, schema versions, aliases, encoded payloads, and message codecs;
- envelope validation, ownership, metadata reservations, and redaction;
- store ordering, expected versions, commit outcomes, iterator behavior, and
  error inspection;
- PostgreSQL schemas, migrations, indexes, transaction and position semantics;
- snapshot and checkpoint encodings;
- queue and Kafka record mappings, keys, headers, settlement, and replay modes;
- metrics, span names, attributes, and cardinality limits;
- module paths, package contracts, minimum Go version, and generated artifacts.

A breaking change to one nested adapter does not require coupling the core
module's release, but every affected module needs its own changelog entry and
directory-prefixed semantic-version tag.

## Release evidence

Before publishing release notes as final, bind them to the immutable release
revision and require the repository's inventory, formatting, tidy, vet, lint,
static analysis, tests, race, exact coverage, mutation, fuzz, vulnerability,
secret, license, SBOM, provenance, docs, API compatibility, generated-code,
integration, clean-consumer, and affected-module gates. Missing or skipped
evidence is a release blocker, not a warning.

Record exact dependency and service-image versions, compatibility-matrix
baseline, migration requirements, known limitations, deprecations, security
impact, rollback boundaries, and raw performance evidence. PostgreSQL, queue,
Kafka, outbox, and telemetry guarantees must remain separate from application
and deployment responsibilities.

## Upgrade workflow

1. read every affected module changelog between the current and target tags;
2. inspect event, wire, schema, error, metric, and repository-contract changes;
3. back up authoritative history and rehearse migrations on a restored copy;
4. run application codecs, upcasters, snapshots, projections, process managers,
   stores, dispatchers, and adapters through their conformance suites;
5. rebuild derived data where the release note requires it;
6. deploy with bounded observation and a rollback or reconciliation procedure;
7. retain old decoders and keys while any live history or backup needs them.

EventSauce conceptual compatibility is tracked separately in the
[versioned matrix](compatibility/eventsauce-3.9.1.md). It does not imply PHP
source compatibility or unspecified wire compatibility.
