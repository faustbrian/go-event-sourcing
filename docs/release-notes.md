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

After v2.0.0 is publicly available, require
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

Independent nested modules retain their existing source locations, module
paths, releases, and core-v1 requirements in this root-only transition. They
are not core-v2 adapters. When a maintained module adopts core-v2 types, its
own major suffix follows its directory: for example,
`github.com/faustbrian/go-event-sourcing/postgres/v2`, not
`github.com/faustbrian/go-event-sourcing/v2/postgres`. Publish core-v2 first,
then PostgreSQL-v2 before an outbox-v2 adapter that depends on both. Queue,
Kafka, and OpenTelemetry adoption may follow core independently.

The deprecated `adapters/gokafka` and `adapters/gotelemetry` v1 implementations
remain unchanged during their existing support interval. Root-v2 availability
does not retire them or change their published contracts. Existing v1
applications can continue using their selected public releases.

`api/v1.0.0.txt` preserves the exact released root API from commit
`9f7bffbb757620e0e0796d9fb564bf66e3cabf7a`. `api/baseline.txt` is the current
root-major API snapshot; cross-major compatibility is an adoption decision,
not an assertion that the two module identities are interchangeable.

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
