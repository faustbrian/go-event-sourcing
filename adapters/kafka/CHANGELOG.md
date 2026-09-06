# Changelog

All notable changes to this module are documented here.

## [Unreleased]

## [1.0.1] - 2026-09-06

### Changed

- Clarify that migration from `adapters/gokafka` changes package-qualified
  concrete type, reflection, and sentinel-error identities even though the
  successor preserves the wire and runtime behavior.
- Supersede the immutable v1.0.0 release with the complete signed release
  artifact set without changing the module's API or runtime behavior.

## [1.0.0] - 2026-09-05

### Added

- Add the target-oriented `adapters/kafka` successor with the complete public
  API, wire format, error classification, bounds, ownership, concurrency, and
  broker behavior of the released `adapters/gokafka` v1 contract.
- Retain the canonical wire version and exact diagnostic categories so moving
  imports does not require broker, record, or failure-policy migration.

### Migration

- New code should import
  `github.com/faustbrian/go-event-sourcing/adapters/kafka`. Existing
  `adapters/gokafka` users may move after this successor is released by changing
  the import path and renaming the qualifier to `kafka`, or by aliasing the new
  import as `gokafka` to preserve existing selectors.

[Unreleased]: https://github.com/faustbrian/go-event-sourcing/compare/adapters%2Fkafka%2Fv1.0.1...HEAD
[1.0.1]: https://github.com/faustbrian/go-event-sourcing/releases/tag/adapters%2Fkafka%2Fv1.0.1
[1.0.0]: https://github.com/faustbrian/go-event-sourcing/releases/tag/adapters%2Fkafka%2Fv1.0.0
