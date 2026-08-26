# Event sourcing OpenTelemetry adapter

`gotelemetry` is the independently versioned observability boundary for
`event-sourcing`. The event-sourcing core does not import OpenTelemetry or the
repository `telemetry` module.

This adapter instruments synchronous dispatch and consumer handling and adds
bounded Kafka context propagation plus event-store append, stream-read, and
global-read instrumentation. Snapshot-store instrumentation observes explicit
load, refresh, and deletion without exposing derived state or aggregate
identity. Projection-runner instrumentation observes bounded replay progress,
poison skips, durable checkpoint position, and terminal probes.

## Install

```sh
go get github.com/faustbrian/go-event-sourcing/adapters/gotelemetry@v1
```

## Quick start

```go
instrumentation, err := gotelemetry.New(runtime)
if err != nil {
	return err
}

dispatcher, err := instrumentation.WrapDispatcher(baseDispatcher)
if err != nil {
	return err
}

handler, err := instrumentation.WrapConsumer(projectEvent)
if err != nil {
	return err
}
consumer, err := eventsourcing.NewConsumer("account_projection", handler)
if err != nil {
	return err
}

kafkaPublisher, err := instrumentation.WrapKafkaPublisher(
	producer,
	gotelemetry.KafkaPropagationConfig{},
)
if err != nil {
	return err
}

kafkaHandler, err := instrumentation.WrapKafkaHandler(
	recordHandler,
	gotelemetry.KafkaPropagationConfig{},
)
if err != nil {
	return err
}

store, err := instrumentation.WrapEventStore(baseStore)
if err != nil {
	return err
}

globalReader, err := instrumentation.WrapGlobalReader(baseGlobalReader)
if err != nil {
	return err
}

snapshotStore, err := instrumentation.WrapSnapshotStore(baseSnapshotStore)
if err != nil {
	return err
}

projectionRunner, err := instrumentation.WrapProjectionRunner(
	"account-summary",
	baseProjectionRunner,
)
if err != nil {
	return err
}

processManager, err := gotelemetry.WrapProcessManager(
	instrumentation,
	"welcome-email",
	baseProcessManager,
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

- [Documentation index](docs/README.md)
- [Complete technical guide](docs/reference.md)
- [Go API reference](https://pkg.go.dev/github.com/faustbrian/go-event-sourcing/adapters/gotelemetry)
- [Parent package documentation](../../docs/README.md)

## Compatibility and support

This module follows Semantic Versioning. Report vulnerabilities through the
[parent security policy](../../SECURITY.md).

## License

MIT. See [LICENSE](LICENSE).
