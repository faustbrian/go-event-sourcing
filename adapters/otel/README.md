# Event sourcing OpenTelemetry adapter

`otel` is the independently versioned observability boundary for
`event-sourcing`. The event-sourcing core does not import OpenTelemetry or the
repository `telemetry` module.

This adapter instruments synchronous dispatch and consumer handling and adds
bounded Kafka context propagation plus event-store append, stream-read, and
global-read instrumentation. Snapshot-store instrumentation observes explicit
load, refresh, and deletion without exposing derived state or aggregate
identity. Projection-runner instrumentation observes bounded replay progress,
poison skips, durable checkpoint position, and terminal probes.

## Lifecycle and Go support

This is the stable, supported OpenTelemetry adapter. It requires Go 1.26.6 and
is tested with Go 1.26.6.

## Install

```sh
go get github.com/faustbrian/go-event-sourcing/adapters/otel@v1
```

## Quick start

Create an application-owned runtime from standard OpenTelemetry providers and
pass it to `otel.New`. The adapter borrows those providers and starts no
goroutines; the application remains responsible for flushing and shutting them
down. The compiling [`ExampleNew`](example_test.go) contains complete imports,
setup, and cleanup.

## Migration from `adapters/gotelemetry`

After this module's first release, change the import path and either rename the
package qualifier from `gotelemetry` to `otel` or explicitly alias the new
import as `gotelemetry`. The public API, signal names, attributes,
instrumentation scope, semantic convention, propagation, privacy bounds, and
provider ownership are unchanged. Release this successor before updating the
deprecated facade to depend on it.

## Guarantees and limitations

The [complete guide](docs/reference.md) defines ownership, failure semantics,
bounds, concurrency, security, and unsupported behavior. Do not infer
additional guarantees beyond the documented module boundary.

## Documentation

This module is listed in the versioned [Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/README.md)
and follows its [persistence and durability family guidance](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/design-language.md#package-families-and-selection).

- [Documentation index](docs/README.md)
- [Complete technical guide](docs/reference.md)
- [Go API reference](https://pkg.go.dev/github.com/faustbrian/go-event-sourcing/adapters/otel)
- [Parent package documentation](../../docs/README.md)

## Compatibility and support

This module follows Semantic Versioning. Report vulnerabilities through the
[parent security policy](../../SECURITY.md).

## License

MIT. See [LICENSE](LICENSE).
