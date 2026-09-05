# Event sourcing OpenTelemetry adapter

`gotelemetry` is the deprecated compatibility path for the preferred
[`adapters/otel`](https://pkg.go.dev/github.com/faustbrian/go-event-sourcing/adapters/otel)
module. It preserves the complete released v1 API, signal names, semantic
convention, and OpenTelemetry instrumentation scope while applications
migrate.

This adapter instruments synchronous dispatch and consumer handling and adds
bounded Kafka context propagation plus event-store append, stream-read, and
global-read instrumentation. Snapshot-store instrumentation observes explicit
load, refresh, and deletion without exposing derived state or aggregate
identity. Projection-runner instrumentation observes bounded replay progress,
poison skips, durable checkpoint position, and terminal probes.

## Lifecycle and Go support

This stable v1 module is deprecated but supported for the compatibility
interval below. It requires Go 1.26.6 and is tested with Go 1.26.6.

## Install

New code should install the v1 successor:

```sh
go get github.com/faustbrian/go-event-sourcing/adapters/otel@v1
```

The released path remains supported for the longer of 180 days and two
published stable minor releases after the successor is publicly consumable.
Removal requires an authorized next-major release:

```sh
go get github.com/faustbrian/go-event-sourcing/adapters/gotelemetry@v1
```

## Quick start

The compiling [`ExampleNew`](example_test.go) constructs instrumentation with
caller-owned OpenTelemetry providers and shuts them down explicitly.

## Guarantees and limitations

The [complete guide](docs/reference.md) defines ownership, failure semantics,
bounds, concurrency, security, and unsupported behavior. Do not infer
additional guarantees beyond the documented module boundary.

## Documentation

This module is listed in the versioned [Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/README.md)
and follows its [persistence and durability family guidance](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/design-language.md#package-families-and-selection).

- [Documentation index](docs/README.md)
- [Complete technical guide](docs/reference.md)
- [Go API reference](https://pkg.go.dev/github.com/faustbrian/go-event-sourcing/adapters/gotelemetry)
- [Parent package documentation](../../docs/README.md)

## Compatibility and support

This stable v1 compatibility module follows Semantic Versioning. When migrating
to `adapters/otel`, rename the package qualifier to `otel` or alias the new
import as `gotelemetry`. Signal names, attributes, instrumentation scope,
propagation, errors, and provider ownership do not change. The deprecated module
retains its released concrete type, error, reflection, and `%T` identities
throughout the support interval. Report vulnerabilities through the
[parent security policy](../../SECURITY.md).

## License

MIT. See [LICENSE](LICENSE).
