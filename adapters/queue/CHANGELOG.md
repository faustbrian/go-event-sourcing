# Changelog

All notable changes to this module will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.0.1] - 2026-10-06

### Changed

- Select Moby client v0.6.1 and API v1.56.1 in the fixture dependency
  graph, with go-connections v0.8.1; retain the owned adapter contracts.

- Adopt Testcontainers core v0.44.0 and Moby client v0.5.0 in the
  dependency closure. The selected otelhttp v0.69.0 can require migration
  for applications directly using removed upstream HTTP helpers; see the
  [supplier adoption guide](https://github.com/faustbrian/go-event-sourcing/blob/main/docs/adapter-migration.md#fixture-supplier-adoption).
  The owned queue and settlement contracts remain unchanged.

- Update indirect compression support to v1.20.0 for integration-
  service setup while preserving the owned delivery and settlement
  contracts.

## [2.0.0] - 2026-10-02

### Changed

- Prepare the independent `/adapters/queue/v2` identity for v2.0.0 and adopt
  public core-v2 nominal types. Preserve Go-Queue v1, envelope bytes,
  ordering, acceptance, and settlement behavior.
- Retain the released v1.0.0 API beside the current-major projection.

- Upgrade the OpenTelemetry SDK dependency graph to v1.45.0 to address
  GHSA-8wmf-6v46-5gfg in transitive dependencies.

- Publish schema-v2 cohesion metadata and enforce it through the repository's
  pinned `go-library-tools` v1.4.0 workflow.
- Reconcile the public v1.0.0 event-sourcing and queue dependency checksums
  with SumDB.

### Documentation

- Move detailed module guidance behind a concise README and documentation index.
- Use human-oriented section names and package-owned documentation links.
- Link the module to the immutable v1.4.0 ecosystem index and
  persistence-and-durability family guidance.

## [1.0.0] - 2026-08-25

### Documentation

- Link the package README to package-owned documentation.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-event-sourcing/adapters/queue` identity while preserving its documented API and behavior.
- Rename the unpublished module from `adapters/goqueue` to `adapters/queue`
  and its Go package to `eventqueue` before v1.
- Require owned sibling modules at local `v0.0.0`; clean external consumers
  pin each module to an exact main pseudo-version.

### Added

- Add a bounded canonical JSON codec that preserves complete event delivery
  identity for compatible queue backends without importing queue behavior into
  the event-sourcing core.
- Add ordered live-only queue dispatch and task handling with separately named
  replay entry points, exact partial-progress errors, panic containment,
  immutable job options, and queue-owned settlement.
- Prove successful and failed delivery through the repository queue and
  in-memory worker while retaining backend-specific guarantee boundaries.
- Prove complete delivery retention and post-handler acknowledgement through
  digest-pinned Valkey Streams 9.1.0 after the producer worker is closed.
- Expose whether the stopping delivery was not attempted or has unknown queue
  acceptance so retry decisions retain duplicate risk.
- Prove hostile wire classes, concurrent shared use, callback and byte
  ownership, accepted-then-failed publication, duplicate retry, blocked
  enqueue cancellation, and application-owned idempotency.
- Add digest-pinned Valkey Streams process-death, disconnection, shutdown,
  ordering-identity, and dead-letter failure recovery evidence, plus separate
  adapter-overhead and durable-publication benchmarks.

### Fixed

- Derive stable event identity for queue failure and dead-letter records,
  preflight the complete first-party queue message bound, reject typed-nil
  queues and unencodable job policy, and copy correlation and trace-context
  metadata across ownership boundaries.
- Redact wrapped input, consumer, panic, backend, and credential diagnostics
  from every common Go error format, including Go-syntax formatting.

[Unreleased]: https://github.com/faustbrian/go-event-sourcing/compare/adapters%2Fqueue%2Fv1.0.0...HEAD
[1.0.0]: https://github.com/faustbrian/go-event-sourcing/releases/tag/adapters%2Fqueue%2Fv1.0.0
