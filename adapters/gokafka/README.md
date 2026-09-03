# Event sourcing Kafka adapter

`gokafka` is the dedicated Kafka boundary for event-sourcing messages. It uses
the repository's franz-go-backed `kafka` module so acknowledgements, producer
idempotence, retries, consumer groups, cooperative rebalancing, offset
settlement, TLS, SASL, and broker operations remain observable Kafka
semantics.

The stable record codec, synchronous direct dispatcher, consumer-group record
handler, explicit poison/retry policy, and first-party synchronous dead-letter
publisher are implemented. Real-broker compatibility covers the complete live
and dead-letter delivery paths.

## Install

```sh
go get github.com/faustbrian/go-event-sourcing/adapters/gokafka@v1
```

## Quick start

```go
codec, err := gokafka.NewRecordCodec(gokafka.RecordCodecConfig{
	Resolver:      gokafka.FixedTopic("accounts.events.v1"),
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

## Guarantees and limitations

The [complete guide](docs/reference.md) defines ownership, failure semantics,
bounds, concurrency, security, and unsupported behavior. Do not infer
additional guarantees beyond the documented module boundary.

## Documentation

This module follows the versioned [Golib ecosystem design language](https://github.com/faustbrian/go-library-tools/blob/v1.3.0/docs/ecosystem/design-language.md)
and is indexed in the [persistence and durability family](https://github.com/faustbrian/go-library-tools/blob/v1.3.0/docs/ecosystem/README.md).

- [Documentation index](docs/README.md)
- [Complete technical guide](docs/reference.md)
- [Go API reference](https://pkg.go.dev/github.com/faustbrian/go-event-sourcing/adapters/gokafka)
- [Parent package documentation](../../docs/README.md)

## Compatibility and support

This module follows Semantic Versioning. Report vulnerabilities through the
[parent security policy](../../SECURITY.md).

## License

MIT. See [LICENSE](LICENSE).
