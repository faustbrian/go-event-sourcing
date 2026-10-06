# Changelog

All notable changes to this module are documented here.

## [Unreleased]

### Changed

- Select Moby client v0.6.1 and API v1.56.1 in the fixture dependency
  graph; retain the owned adapter contracts.

- Adopt Testcontainers core v0.44.0 and Moby client v0.5.0 in the
  dependency closure. The selected otelhttp v0.69.0 can require migration
  for applications directly using removed upstream HTTP helpers; see the
  [supplier adoption guide](https://github.com/faustbrian/go-event-sourcing/blob/main/docs/adapter-migration.md#fixture-supplier-adoption).
  Kafka fixtures use v0.44.0; owned adapter contracts remain unchanged.
- Adopt franz-go v1.22.1, kmsg v1.14.0, compression v1.20.0 and LZ4
  v4.1.30 in the selected Kafka dependency graph. Keep the owned wire,
  settlement and telemetry contracts and explicit client bounds.
  Upstream clients now reject deleted-and-recreated topics and require
  explicit BalanceRacks to opt into rack-aware group assignment;
  review these deployment cases before upgrading.
- Prepare the independent `/adapters/otel/v2` identity for v2.0.0 and adopt
  public core-v2 nominal types. Preserve dependencies, signal names, the
  unsuffixed instrumentation scope, and deprecated gotelemetry-v1 support.
- Retain the released v1.0.2 API beside the current-major projection.

### Fixed

- Upgrade the OpenTelemetry SDK dependency graph to v1.45.0 to address
  GHSA-8wmf-6v46-5gfg, a conditional exporter-configuration disclosure through
  verbose internal SDK logs.

## [1.0.2] - 2026-09-06

### Changed

- Supersede the immutable v1.0.1 release to correct the `SHA256SUMS` SSH
  signature namespace from `file` to the canonical `golib-release` namespace.
  The module API and runtime behavior are unchanged.

## [1.0.1] - 2026-09-06

### Changed

- Clarify that migration from `adapters/gotelemetry` changes
  package-qualified concrete type, reflection, and sentinel-error identities
  even though the successor preserves its signals and runtime behavior.
- Supersede the immutable v1.0.0 release with the complete signed release
  artifact set without changing the module's API or runtime behavior.

## [1.0.0] - 2026-09-05

### Added

- Add the target-oriented `adapters/otel` successor with the complete public
  API, telemetry signals, privacy bounds, propagation behavior, error
  classification, and lifecycle semantics of the released
  `adapters/gotelemetry` v1 contract.
- Preserve the existing OpenTelemetry instrumentation scope and semantic
  convention so moving imports does not split traces, metrics, or dashboards.

### Migration

- New code should import
  `github.com/faustbrian/go-event-sourcing/adapters/otel`. Existing
  `adapters/gotelemetry` users may move after this successor is released by
  changing the import path and renaming the qualifier to `otel`, or by aliasing
  the new import as `gotelemetry` to preserve existing selectors.

[Unreleased]: https://github.com/faustbrian/go-event-sourcing/compare/adapters%2Fotel%2Fv1.0.2...HEAD
[1.0.2]: https://github.com/faustbrian/go-event-sourcing/releases/tag/adapters%2Fotel%2Fv1.0.2
[1.0.1]: https://github.com/faustbrian/go-event-sourcing/releases/tag/adapters%2Fotel%2Fv1.0.1
[1.0.0]: https://github.com/faustbrian/go-event-sourcing/releases/tag/adapters%2Fotel%2Fv1.0.0
