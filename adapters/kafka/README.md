# Event sourcing Kafka adapter

`kafka` is the dedicated Kafka boundary for event-sourcing messages. It uses
the repository's franz-go-backed `kafka` module so acknowledgements, producer
idempotence, retries, consumer groups, cooperative rebalancing, offset
settlement, TLS, SASL, and broker operations remain observable Kafka
semantics.

The stable record codec, synchronous direct dispatcher, consumer-group record
handler, explicit poison/retry policy, and first-party synchronous dead-letter
publisher are implemented. Real-broker compatibility covers the complete live
and dead-letter delivery paths.

## Lifecycle and Go support

This is the stable, supported Kafka adapter. It requires Go 1.26.6 and is
tested with Go 1.26.6.

## Install

```sh
go get github.com/faustbrian/go-event-sourcing/adapters/kafka@v1
```

## Quick start

```go
codec, err := kafka.NewRecordCodec(kafka.RecordCodecConfig{
	Resolver:      kafka.FixedTopic("accounts.events.v1"),
	AllowedTopics: []string{"accounts.events.v1"},
})
if err != nil {
	return err
}

record, err := codec.Encode(delivery)
if err != nil {
	return err
}
err = producer.Publish(ctx, record)
```

The compiling examples in this module contain complete imports and setup.

## Migration from `adapters/gokafka`

Change the import path and either rename the package qualifier from `gokafka`
to `kafka` or explicitly alias the new import as `gokafka`. The API shape, wire
version, headers, ordering, bounds, error categories, and ownership contracts
are unchanged. Package-qualified concrete type, reflection, and sentinel error
identities change with the import path. The deprecated module remains an
independent implementation during its support interval.

## Guarantees and limitations

The [complete guide](docs/reference.md) defines ownership, failure semantics,
bounds, concurrency, security, and unsupported behavior. Do not infer
additional guarantees beyond the documented module boundary.

## Documentation

This module is listed in the versioned [Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/README.md)
and follows its [persistence and durability family guidance](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/design-language.md#package-families-and-selection).

- [Documentation index](docs/README.md)
- [Complete technical guide](docs/reference.md)
- [Go API reference](https://pkg.go.dev/github.com/faustbrian/go-event-sourcing/adapters/kafka)
- [Parent package documentation](../../docs/README.md)

## Compatibility and support

This module follows Semantic Versioning. Report vulnerabilities through the
[parent security policy](../../SECURITY.md).

## License

MIT. See [LICENSE](LICENSE).
