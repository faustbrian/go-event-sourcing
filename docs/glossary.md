# Glossary

**Aggregate**<br>
A consistency boundary that owns domain invariants and changes state by
recording events.

**Aggregate root ID**<br>
The application-defined stable identity of one aggregate stream. It need not be
a UUID.

**Aggregate type**<br>
The stable application-defined category paired with an aggregate root ID to
form a stream identity.

**Causation ID**<br>
The message ID of the immediate input that caused another message.

**Checkpoint**<br>
A durable global position through which a projection has successfully
processed or explicitly skipped messages.

**Commit ambiguity**<br>
A failure where the caller cannot prove whether a transaction committed. It
requires reconciliation, not blind retry.

**Correlation ID**<br>
An identifier used to associate messages in one business interaction without
changing their aggregate identity.

**Delivery**<br>
A persisted event message plus an explicit live or replay mode supplied to a
consumer.

**Dispatcher**<br>
The replaceable responsibility that delivers persisted messages to consumers.

**Event**<br>
An immutable statement of a domain transition that has already occurred.

**Event name**<br>
The explicit stable persisted identity of an event. It is independent of Go
package paths and type names.

**Event schema version**<br>
The positive version describing the persisted payload contract for an event
name.

**Event store**<br>
The replaceable responsibility that atomically appends and reads immutable
messages with expected-version semantics.

**Expected version**<br>
The concurrency precondition for an append: new, existing, exact, or explicit
any-version mode.

**Global position**<br>
An optional one-based store-wide ordering position assigned by capable stores.

**Live delivery**<br>
Delivery of a message after its normal durable append path.

**Message**<br>
The immutable persisted envelope containing event data, stream identity,
versions, time, metadata, and causal identifiers.

**Message ID**<br>
The stable unique identity used for duplicate detection and reconciliation.

**Optimistic concurrency**<br>
Rejecting an append when the actual stream version no longer matches its
declared precondition.

**Outbox**<br>
Derived durable publication work committed atomically with application or event
data and relayed at least once after commit.

**Pending message**<br>
A validated event envelope before the store assigns stream and optional global
positions.

**Process manager**<br>
A consumer that reacts to events by planning explicit commands or messages.
Effects and persistence remain application-owned.

**Projection**<br>
Derived state built by consuming ordered messages, commonly with a durable
checkpoint. Event sourcing does not require projections or CQRS.

**Replay delivery**<br>
Explicit historical delivery used for rebuild or analysis. Safe compositions
isolate external side effects and process managers by default.

**Snapshot**<br>
Replaceable derived aggregate state at a known aggregate version. Event history
remains authoritative.

**Stream**<br>
The ordered event history for one aggregate type and aggregate root ID.

**Stream version**<br>
The one-based position of a message within its aggregate stream.

**Upcaster**<br>
A deterministic read-boundary transformation from stored event identity,
payload, and metadata to newer logical events without rewriting history.
