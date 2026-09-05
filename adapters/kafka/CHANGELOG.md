# Changelog

All notable changes to this module are documented here.

## Unreleased

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
