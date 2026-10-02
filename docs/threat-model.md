# Event-sourcing threat model

Model version: 1. Reviewed: 2026-10-02. Repository maintainers own this model
and the first-party implementation controls. This is a boundary model, not a
deployment certification or evidence that every security gate has passed.

## Scope and assets

The [module manifest](../modules.json) is the inventory authority. Production
boundaries include the root core, memory stores, process managers, projections,
snapshots, PostgreSQL, and the outbox, queue, Kafka, and OpenTelemetry adapters.
The deprecated `gokafka` and `gotelemetry` modules retain the corresponding
Kafka and telemetry implementation boundaries; deprecation does not remove
their data exposure. `eventtest` is test support and `benchmarks/competitors`
is non-production comparison code, not a persistence or security provider.

Protect authoritative event bytes, ordering and message identities, aggregate
state, derived snapshots and projections, checkpoints, tenant routing data,
outbox/broker copies, and diagnostic privacy. The owned CloudEvents
`adapters/golib` consumer crosses the root envelope boundary; nested adapters
consume root contracts without acquiring application authorization authority.

Stored database rows and broker envelopes are untrusted data. Application
callbacks, custom stores/codecs, database and broker clients, entropy readers,
and telemetry providers are trusted extension code with explicit ownership
and bounded-work responsibilities. A compromised extension is not sandboxed.

## Controls and acceptance boundaries

| Threat / entry point | First-party control | Contract and evidence location |
| --- | --- | --- |
| Oversized or mutable envelope input | Constructors validate before copying; payloads are at most 1 MiB, metadata at most 64 entries / 64 KiB, snapshot state at most 8 MiB; identifier and value limits are explicit | [message.go](../message.go), [snapshot.go](../snapshot.go), envelope validation tests |
| Ambiguous JSON or excessive parsing work | Built-in JSON codec rejects invalid UTF-8, duplicate keys, trailing data, depth above 100 and containers above 10,000 entries; transport codecs bound envelopes and validate canonical fields | [codec.go](../codec.go), [serialization guide](serialization.md), adapter codec tests |
| Unbounded schema evolution | Reference upcaster bounds path steps, output segments and total work, checks deterministic progress, and requires an explicit reviewed drop policy | [upcast.go](../upcast.go), [upcast tests](../upcast_test.go) |
| History or checkpoint corruption | Sequential stream/global positions, expected versions, append identity checks, checkpoint compare-and-set and authoritative snapshot-version checks fail closed | [repository.go](../repository.go), [projection runner](../projection/runner.go), [PostgreSQL guide](../postgres/README.md) |
| Accidental replay side effects | Delivery mode is explicit; process managers and transport handlers reject replay by default; every projection batch requires a replay guard | [process managers](process-managers.md), [projection guard](../projection/guard.go), handler and guard tests |
| Diagnostic disclosure | Built-in codec/registration errors omit input identities and content types; PostgreSQL metadata parser errors expose stable categories without printing stored input while retaining inspectable causes | [privacy regressions](../diagnostic_privacy_test.go), [PostgreSQL regressions](../postgres/iterator_internal_test.go), [security guidance](security.md) |
| Forged history | Optional verifying store/reader checks each message before exposure and terminates on rejection, panic or cancellation | [integrity.go](../integrity.go), [security guidance](security.md#integrity-and-recovery) |
| Ambiguous commit or duplicate delivery | Stable prepared IDs and explicit commit outcomes support reconciliation; outbox staging does not claim that a caller-owned transaction committed; checkpoints advance only after handling | [outbox guide](outbox.md), [repository.go](../repository.go), [projection guide](projections.md) |
| Transport ambiguity and tombstones | Kafka codec checks topic allowlists, routing, reserved headers and timestamps; an empty value is not a valid event payload and does not become an implicit deletion event | [Kafka codec](../adapters/kafka/codec.go), [Kafka guide](kafka.md) |

These controls do not provide authentication, authorization, encryption,
exactly-once application effects, physical erasure, or fleet-wide quotas.
`Tenant` and `Partition` are routing values, not access checks. A replay guard
does not authorize other reads or administrative actions.

## Owned residual risks

| Residual and rationale | Owner and mitigation | Review trigger |
| --- | --- | --- |
| Bounded records do not bound retained history or aggregate replay duration; memory maps and external storage grow across calls | Application/deployment owner: tenant quotas, admission control, caller deadlines, retention and memory-store lifetime limits | New externally exposed write/read route, new retention policy, larger limits or storage backend |
| A synchronous trusted callback, custom codec/upcaster or entropy reader can block, allocate, panic or disclose its own data; context cannot preempt arbitrary Go code | Application extension owner: bounded deterministic callbacks, concurrency-safe providers and safe diagnostics; use context-aware extensions where offered | Replacement callback/provider, new callback output shape or cancellation expectation |
| Driver/broker buffering occurs before this package can validate a row or record; cancellation and safe cleanup depend on the client contract | Deployment/client owner: bounded fetch and pool settings, database constraints, deadlines, least privilege and client updates | Driver/broker upgrade, pool/fetch changes, schema or connection-role change |
| Post-commit delivery and checkpoint failures may repeat already-applied effects; rebuild reset and checkpoint reset are not one atomic operation | Application owner: idempotent consumers, stable-ID reconciliation, explicit poison policy, paused/drained/serialized rebuild controls | New side effect, retry policy, outbox integration, rebuild or repair workflow |
| Routing metadata and topic allowlists do not authenticate a caller; verifiers do not choose a signing scheme or manage keys | Application security owner: authorize every read/write/replay/admin path, choose canonical authenticated history and key lifecycle | New tenant boundary, signing format, verifier, key rotation or replay tool |
| Copies in dead letters, outboxes, snapshots, projections, logs and backups outlive a deletion event or Kafka tombstone | Data/deployment owner: inventory copies, minimize sensitive fields, encrypted transport/storage, explicit retention/erasure procedures | New export/copy, retention requirement, restore or erasure operation |
| Inspecting an unwrapped parser/driver/application cause or explicitly formatting an envelope can expose application data | Application observability owner: log safe outer categories, restrict cause inspection and envelope formatting to authorized diagnostics | Error/logging middleware, tracing exporter or support-bundle change |

## Review and reporting

The nested module graphs pin OpenTelemetry SDK v1.45.0 or later for
[GHSA-8wmf-6v46-5gfg](https://github.com/open-telemetry/opentelemetry-go/security/advisories/GHSA-8wmf-6v46-5gfg).
The affected older SDK versions can log exporter endpoint configuration when
an application enables verbose internal diagnostics while constructing its
provider. The default logger does not emit that event. This package accepts
application-owned providers and does not configure exporter credentials or the
global logger; applications must also upgrade their own provider/exporter
graphs and review existing diagnostic logs where that configuration was used.

Reassess the affected row when a public contract, parser, limit, schema, client,
callback boundary or consumer changes. Preserve stable error categories and
review deprecated variants when a shared implementation changes. Runtime,
external-service, release and consumer evidence remain separate from source
review and documentation checks.

Use [private vulnerability reporting](https://github.com/faustbrian/go-event-sourcing/security/advisories/new)
for suspected vulnerabilities; see [SECURITY.md](../SECURITY.md) for support
and disclosure policy. Do not submit credentials, production payloads or
weaponized demonstrations in public issues.
