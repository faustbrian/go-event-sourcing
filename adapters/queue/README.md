# Event sourcing queue adapter

`eventqueue` maps complete event-sourcing deliveries to the first-party
`github.com/faustbrian/go-queue` contract. The event-sourcing core does
not import queue, and this adapter owns no worker, broker connection, retry
clock, dead-letter store, or business-idempotency state.

The adapter provides a bounded canonical payload, synchronous publication in
input order, explicit enqueue ambiguity, and queue-owned settlement. It does
not claim exactly-once delivery or broker-neutral durability and ordering.

## Install

```sh
go get github.com/faustbrian/go-event-sourcing/adapters/queue@v1
```

## Quick start

```go
codec, err := eventqueue.NewCodec(eventqueue.CodecConfig{})
if err != nil {
	return err
}

handler, err := eventqueue.NewTaskHandler(codec, applyDelivery)
if err != nil {
	return err
}

workerQueue := queue.NewPool(
	1,
	queue.WithFn(handler.Handle),
	queue.WithLogger(queue.NewEmptyLogger()),
)
workerQueue.Start()
defer workerQueue.Release()

dispatcher, err := eventqueue.NewDispatcher(eventqueue.DispatcherConfig{
	Queue: workerQueue,
	Codec: codec,
	Job: job.AllowOption{
		Metadata: &job.Metadata{
			RetryPolicy: "projection-v1",
			HandlerType: "account-projector",
		},
	},
})
if err != nil {
	return err
}

return dispatcher.Dispatch(ctx, deliveries)
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
- [Go API reference](https://pkg.go.dev/github.com/faustbrian/go-event-sourcing/adapters/queue)
- [Parent package documentation](../../docs/README.md)

## Compatibility and support

This module follows Semantic Versioning. Report vulnerabilities through the
[parent security policy](../../SECURITY.md).

## License

MIT. See [LICENSE](LICENSE).
