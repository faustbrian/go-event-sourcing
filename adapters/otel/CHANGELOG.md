# Changelog

All notable changes to this module are documented here.

## [Unreleased]

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

[Unreleased]: https://github.com/faustbrian/go-event-sourcing/compare/adapters%2Fotel%2Fv1.0.1...HEAD
[1.0.1]: https://github.com/faustbrian/go-event-sourcing/releases/tag/adapters%2Fotel%2Fv1.0.1
[1.0.0]: https://github.com/faustbrian/go-event-sourcing/releases/tag/adapters%2Fotel%2Fv1.0.0
