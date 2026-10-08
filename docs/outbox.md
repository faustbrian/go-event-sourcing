# Transactional outbox integration

The core event store and outbox remain independently usable and releasable.
Neither imports the other. The optional
`github.com/faustbrian/go-event-sourcing/adapters/outbox/v3` nested module
is the only component that depends on both public contracts.

The v3 adapter is prepared here and pending publication. It adopts public
Transactional Outbox v2 while retaining Event Sourcing core/PostgreSQL v2.
See its [migration guide](../adapters/outbox/docs/reference.md#v3-migration).

Use its `Stager` with an already caller-owned `pgx.Tx` when event rows and
publishable outbox envelopes must commit together. The application prepares the
aggregate save, stages both batches, commits, and only then confirms and
dispatches the aggregate. The adapter never owns transaction completion,
publishes before commit, or enqueues records during replay.

The complete API, envelope mapping, crash matrix, limits, recovery procedure,
and examples are documented in the
[eventoutbox adapter guide](../adapters/outbox/README.md).
Applications replacing either side should follow the
[custom outbox boundary](custom-outbox.md) without coupling the two cores.

The production path is:

1. append event messages and outbox envelopes in one PostgreSQL transaction;
2. commit PostgreSQL before any external publication;
3. claim committed envelopes through the independently operated outbox relay;
4. publish through the outbox-owned Kafka publisher; and
5. settle Kafka consumption only after successful idempotent handling.

This is at-least-once delivery. PostgreSQL and Kafka do not share a transaction.
Producer idempotence and Kafka transactions do not change that boundary.
