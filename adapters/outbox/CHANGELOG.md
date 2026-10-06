# Changelog

All notable changes to this module are documented here.

## Unreleased

### Changed

- Select Moby client v0.6.1 and API v1.56.1 in the fixture dependency
  graph, with go-connections v0.8.1; retain the owned adapter contracts.

- Adopt Testcontainers core v0.44.0 and Moby client v0.5.0 in the
  dependency closure. The selected otelhttp v0.69.0 can require migration
  for applications directly using removed upstream HTTP helpers; see the
  [supplier adoption guide](https://github.com/faustbrian/go-event-sourcing/blob/main/docs/adapter-migration.md#fixture-supplier-adoption).
  PostgreSQL fixtures retain their v0.43.0 module over core v0.44.0.
- Update indirect compression support to v1.20.0 for integration-
  service setup while preserving the owned atomic staging and outbox
  contracts.
- Prepare the independent `/adapters/outbox/v2` module for v2.0.0, adopting
  public core-v2 and PostgreSQL-v2 together. Public nominal types change;
  callers must migrate all three identities together. Keep Transactional
  Outbox v1 and all staging, codec, SQL, and transaction behavior unchanged.
- Preserve the released v1 API snapshot alongside the current-major API.

- Upgrade the OpenTelemetry SDK dependency graph to v1.45.0 to address
  GHSA-8wmf-6v46-5gfg in transitive dependencies.

- Publish schema-v2 cohesion metadata and enforce it through the repository's
  pinned `go-library-tools` v1.4.0 workflow.
- Reconcile the public v1.0.0 event-sourcing and transactional-outbox
  dependency checksums with SumDB.

### Documentation

- Move detailed module guidance behind a concise README and documentation index.
- Use human-oriented section names and package-owned documentation links.
- Link the module to the immutable v1.4.0 ecosystem index and
  persistence-and-durability family guidance.

## 1.0.0 - 2026-08-25

### Documentation

- Link the package README to package-owned documentation.

### Removed

- Remove the pre-release committed `Store`; callers now retain exclusive
  transaction, commit-ambiguity, aggregate acknowledgement, and dispatch
  ownership through `Stager`.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-event-sourcing/adapters/outbox` identity while preserving its documented API and behavior.
- Upgrade `golang.org/x/text` to v0.41.0 so the dependency graph no longer
  contains GO-2026-5970.
- Rename the unpublished module from `adapters/gooutbox` to
  `adapters/outbox` and its Go package to `eventoutbox` before v1.
- Require owned sibling modules at local `v0.0.0`; clean external consumers
  pin each module to an exact main pseudo-version.
- Isolate every staging attempt in a PostgreSQL savepoint so mapping, database,
  cancellation, and savepoint failures cannot leave a one-sided batch
  committable by the caller while outer rollback remains caller-owned.
- Resolve custom topics from the new pre-persistence `TopicMessage` contract
  before acquiring stream locks. Custom resolvers must change their parameter
  from `eventsourcing.Message` to `eventoutbox.TopicMessage`; store-assigned stream
  versions and global positions are no longer available for routing.

### Added

- failure-injection, concurrency, ambiguity-recovery, process-death, relay
  duplication, hostile-codec, and realistic batch-performance evidence
- a digest-pinned PostgreSQL 18.4 fixture for same-transaction integration
  evidence
- same-transaction event and outbox staging directly from a prepared core
  aggregate save plan
- real PostgreSQL benchmark isolating same-transaction outbox append overhead
- complete committed-store-to-relay composition guidance and real PostgreSQL
  evidence for transient retry, delivery, and replay isolation
- atomic committed event-and-outbox storage through `Store`
- caller-owned same-transaction staging through `Stager`
- deterministic bounded event-message envelope mapping and validation
- explicit commit-outcome, rollback, replay, and at-least-once semantics
