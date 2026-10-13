// Package eventsourcing provides storage-independent contracts and explicit
// composition for event-sourced domain models. Start with NewRepository and
// explicit codecs, stores, clocks, identifiers and dispatchers. Core and memory
// composition use the standard library; PostgreSQL and transport/telemetry
// integrations are independently released optional modules.
//
// Persistence and post-commit delivery are separate outcomes. Unknown commits
// require reconciliation, not blind retry. Replay preserves stored source
// positions even when upcasting splits or drops logical events. Application
// values/callbacks retain caller-owned immutability, concurrency and blocking
// obligations; the core owns no service loop, transaction or shutdown.
//
// Error strings omit selected diagnostics but wrapped causes/public fields can
// retain sensitive data. Identities and tenant metadata are not authentication
// or authorization, and message integrity requires an explicit verifier.
package eventsourcing
