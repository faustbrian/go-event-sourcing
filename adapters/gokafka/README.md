# Event sourcing Kafka adapter

`gokafka` is the deprecated compatibility path for the preferred
[`adapters/kafka`](https://pkg.go.dev/github.com/faustbrian/go-event-sourcing/adapters/kafka)
module. It preserves the complete released v1 API, wire format, error
categories, and synchronous Kafka behavior while applications migrate.

The stable record codec, synchronous direct dispatcher, consumer-group record
handler, explicit poison/retry policy, and first-party synchronous dead-letter
publisher are implemented. Real-broker compatibility covers the complete live
and dead-letter delivery paths.

## Lifecycle and Go support

This stable v1 module is deprecated but supported for the compatibility
interval below. It requires Go 1.26.6 and is tested with Go 1.26.6.

## Install

New code should install the successor after its first release:

```sh
go get github.com/faustbrian/go-event-sourcing/adapters/kafka@v1
```

The released path remains supported for the longer of 180 days and two
published stable minor releases after the successor is publicly consumable.
Removal requires an authorized next-major release:

```sh
go get github.com/faustbrian/go-event-sourcing/adapters/gokafka@v1
```

## Quick start

The executable [`ExampleNewRecordCodec`](facade_test.go) compiles representative
construction through the released path. The same file compares stable error
categories with the successor.

## Guarantees and limitations

The [complete guide](docs/reference.md) defines ownership, failure semantics,
bounds, concurrency, security, and unsupported behavior. Do not infer
additional guarantees beyond the documented module boundary.

## Documentation

This module is listed in the versioned [Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/README.md)
and follows its [persistence and durability family guidance](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/design-language.md#package-families-and-selection).

- [Documentation index](docs/README.md)
- [Complete technical guide](docs/reference.md)
- [Go API reference](https://pkg.go.dev/github.com/faustbrian/go-event-sourcing/adapters/gokafka)
- [Parent package documentation](../../docs/README.md)

## Compatibility and support

This stable v1 compatibility module follows Semantic Versioning. When migrating
to `adapters/kafka`, rename the package qualifier to `kafka` or alias the new
import as `gokafka`. Topics, records, acknowledgements, failure policies, error
categories, and runtime ownership do not change. Facade aliases expose the
successor package path through reflection and `%T`; applications must not use
concrete package identity as a behavior or persistence key.
Report vulnerabilities through the [parent security policy](../../SECURITY.md).

## License

MIT. See [LICENSE](LICENSE).
