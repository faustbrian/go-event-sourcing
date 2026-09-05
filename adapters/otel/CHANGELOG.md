# Changelog

All notable changes to this module are documented here.

## Unreleased

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
