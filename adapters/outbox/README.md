# Event sourcing outbox adapter

`eventoutbox` composes the public event-sourcing PostgreSQL writer and the public
outbox PostgreSQL writer. Neither core imports the other. `Stager` writes event
rows and one outbox envelope per event through a savepoint inside an existing
caller-owned `pgx.Tx`. It releases that savepoint only after both batches stage
and rolls it back after any error. The adapter never commits or publishes from
the outer transaction and never claims exactly-once delivery.

## Install

Adapter v3 is prepared in this source and pending publication. After release:

```sh
go get github.com/faustbrian/go-event-sourcing/adapters/outbox/v3@v3.0.0
```

Use Transactional Outbox `/v2` envelopes, limits and PostgreSQL writers with
this adapter. Event Sourcing core `/v2` and PostgreSQL `/v2` stay unchanged;
the selected minimums include their v2.0.1 and v2.0.2 fixes respectively.

## Quick start

```go
limits := eventoutbox.DefaultLimits()
outboxWriter, err := outboxpostgres.NewWriter(
	outboxpostgres.WriterConfig{
		Limits:       limits,
		MaxBatchSize: eventsourcing.MaxAppendMessages,
	},
)
if err != nil {
	return err
}
codec, err := eventoutbox.NewEnvelopeCodec(
	eventoutbox.FixedTopic("account-events"),
	limits,
)
if err != nil {
	return err
}
```

The compiling examples in this module contain complete imports and setup.

## Guarantees and limitations

The [complete guide](docs/reference.md) defines ownership, failure semantics,
bounds, concurrency, security, and unsupported behavior. Do not infer
additional guarantees beyond the documented module boundary.

## Documentation

This module is listed in the versioned [Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/README.md)
and follows its [persistence and durability family guidance](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/design-language.md#package-families-and-selection).

- [Documentation index](docs/README.md)
- [Complete technical guide](docs/reference.md)
- [Go API reference](https://pkg.go.dev/github.com/faustbrian/go-event-sourcing/adapters/outbox/v3)
- [Parent package documentation](../../docs/README.md)

## Compatibility and support

This module follows Semantic Versioning. Report vulnerabilities through the
[parent security policy](../../SECURITY.md).

See [v3 migration](docs/reference.md#v3-migration) before changing imports.
Existing adapter-v2 callers retain the published outbox-v1 type cohort;
selecting producer v2 alone does not migrate their nominal contracts.

## License

MIT. See [LICENSE](LICENSE).
