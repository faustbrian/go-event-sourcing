# Source documentation audit

[Documentation index](README.md)

## Owner, acceptance and selected checks

This Tier A pass owns source comments and directly contradicted current prose at
fresh remote main `21eba2c6ae75c5c1df4c90ad7278156fbbdce55a`. The original
unrelated checkout is preserved. The coordinator owns enforcement, publication
and exact-main CI. No runtime, dependency, SQL, corpus, historical authority or
desired-policy repairs are in scope.

Consumers include root aggregate repositories/codecs/lifecycle and dispatch,
memory and PostgreSQL stores, projections/process managers/snapshots, assurance
helpers, canonical and independently published compatibility Kafka/telemetry
bindings, transactional-outbox and queue adapters, and the competitor harness.
Named risks are commit versus dispatch certainty, append/replay order and split
events, receipts/checkpoint fencing, cancellation and cleanup ownership, mutable
payload/callback borrowing, safe diagnostics, adapter propagation and consistency,
and live-service/benchmark claims that a source audit cannot establish.

Before source writes, tracked inventory and `golib source-docs inventory --json`
show 283 Go files (80 source, 203 tests), nine modules and 14 source directories.
The initial 267 findings are undocumented exported members/methods, separate
from semantic contracts needing investigation. Tagged integration tests,
assurance helpers, fuzz seeds, examples and benchmark sources are included.
The configured eight API baselines use `apidiff`: preserve their frozen bytes
and prove unchanged signatures/AST rather than regenerate prose into them.

After an explicit serialized Go 1.27.2 grant, use three disposable caches for
native formatting, all-file AST/raw-token/literal/directive preservation,
actual package GoDoc and applicable native docs/Examples. Test-only benchmark
namespaces need source/AST documentation, not competitor execution. Select no
database/broker/service, race, fuzz, mutation, benchmark or release campaign.
Static source/local-documentation/specification applicability, immutable pinned
full/changed/staged secrets, frozen-input comparison and ordinary complete-diff
truth/usefulness inspection finish the local commit. Missing comments and
unresolved runtime/rationale questions remain distinct below.

## Source inventory and investigation

Tracked inventory is reproducible with `git ls-files '*.go' '**/go.mod'` and
`golib source-docs inventory --json`, independent of selected build packages.
There is no generated Go or vendor source. All 33 constrained files are
integration-tagged test/benchmark/helper sources, not alternate production
implementations; their directives and bodies stay unchanged. The two embedded
SQL migrations are hand-maintained forward-only inputs with no Down directives.
Their bytes, embed declaration, all dependency/catalog/policy inputs, eight
structural API baselines, frozen wire corpus and historical authority remain
unchanged. No specification decision set or official conformance claim is
configured. `docs/code-generation.md` explicitly records that no generator is
shipped; headers must not invent one.

| Module/source boundary | Contracts investigated |
| --- | --- |
| Root plus memory/projection/processmanager/snapshot/eventtest | Immutable bounded encoded messages versus shallow application values, identity without authorization, callback/context borrowing, per-stream atomic append and unknown commit classification, lifecycle poisoning without state rollback, stored-source split/drop replay, stage/confirm/dispatch ownership, projection handler effects before checkpoint CAS, explicit authorization/poison/rebuild hooks, derived snapshots versus authoritative history, and fixture/callback ownership. |
| postgres | Borrowed pools and single-owner rows/iterators, per-writer transaction permits versus transaction-wide caller serialization, per-stream lock plus global commit-order row, conservative commit reconciliation on a primary, positive bigint bounds, post-scan validation versus pre-transfer limits, forward-only migrations and fresh five-second background rollback attempts. |
| adapters/kafka and adapters/gokafka | Separate v2/core-v2 and compatibility v1/core-v1 identities, shared package wire-v1 mapping, copied exact topic allowlists/limits, synchronous acknowledgement with ambiguous failures, caller broker partitioning/order, returned-nil settlement signals, replay opt-in, owned dead-letter records and cancellation after successful publication. |
| adapters/otel and adapters/gotelemetry | Separate published core bindings, unchanged scope, borrowed providers/instruments/exporter shutdown, isolated telemetry panics but propagated business panics, original caller context authority, bounded traceparent/tracestate copies, borrowed inbound payload/key, operation spans lasting until explicit Close and partial replay progress passed through unchanged. |
| adapters/outbox/v3 | Public core-v2/PostgreSQL-v2/outbox-v2 composition, pre-lock deterministic topic resolution, one savepoint for staged pairs, caller-owned outer commit/rollback/dispatch/reconciliation, cleanup-failure abort attempts and no durable success merely from staging. |
| adapters/queue/v2 | Package-defined bounded canonical wire envelope, copied job/metadata pointers, synchronous producer with no context parameter, unknown acceptance on failed/panicking admission, replay policy, task accessor/consumer panic containment and queue-owned settlement. |
| benchmarks/competitors | Test-only namespace with pinned third-party dependencies, selected equivalent state/version fixtures, declared history lengths, excluded setup versus measured record/apply/replay; no public library or production performance qualification. |

All 203 test/helper files were inventoried by named outcome and inspected at
contract assertions: ownership/validation/privacy, split evolution and lifecycle
generations, commit/dispatch receipts and unknown outcomes, checkpoint winners
and lock cancellation, snapshot compatibility/equivalence, transport mapping and
settlement, hostile provider/context isolation, and optional broker/database
restart/cleanup/allocator/replica cases. The 32 fuzz targets now identify their
actual input/observable boundary. Benchmark source documents workload/setup and
comparison limits; immutable saved corpora and expected outcomes are retained.
There are eight root and two telemetry output-bearing Examples; four Kafka
Examples deliberately lack Output and require an actual broker/operator loop.
Their compile-only status is not execution evidence.

## Findings and source-only questions

The initial 267 member/method documentation omissions are resolved. Canonical
package comments replace duplicate compatibility package headers while keeping
valid Deprecated markers and the current /v2 migration target. Ordinary source
structure is zero findings, not certification of semantic completion.

Directly contradicted maintained prose/comments were corrected for no-I/O save
preparation, unconditional dispatcher shared-use safety, per-topic keying versus
broker ordering, process-local RAM durability, reversible migration wording,
and iterator/context ownership. Most substantive contracts already had useful
source comments; the pass retains those instead of rewriting every declaration.

Concrete source boundaries, not repaired runtime defects: memory snapshot
methods check cancellation only before waiting for their mutex; copied old
SaveResult values can be dispatched again; application JSON marshalers can
allocate/block/panic before encoded size validation; PostgreSQL validation occurs
after driver scanning. Deferred rollback attempts mostly ignore cleanup failure.
Outbox savepoint cleanup failure triggers a bounded outer-abort attempt and
returns an error; caller rollback remains mandatory, not a source proof of
remote transaction cleanup. These statements do not replace stronger desired
policy or qualify every fault/interoperability scenario.

Checkpoint CAS protects progress, not earlier handler effects or a durable worker
lease. SDK exporters, callback budgets, real identity authorization, SQL/broker
durability and distributed side-effect idempotency remain application/deployment
responsibilities. No rationale for alternate behavior is invented. Legacy
adapters retain public v1 types rather than silently acquiring current root-v2
contracts, and historical release/order/compatibility authority remains frozen.

## Local proof and limits

The explicitly serialized Go 1.27.2 session passed formatting and exact runtime
AST/signature/raw-token/literal/implicit-semicolon comparison for all 283 files;
no token normalization was needed. All directives, analyzer annotations and
Example Output comments remain exact. All 1,450 comment groups parsed/rendered
as GoDoc HTML. Thirteen production packages rendered with `go doc -all .`;
the competitor test-only namespace rendered through the standard-library AST
and GoDoc APIs without compiling or running its harness. Final member wording
refinements refreshed all-file preservation and affected root/eventtest renders;
unchanged native docs and Examples retain their original results.

The six configured module `make docs` operations passed. Root and both telemetry
bindings passed `go test -run '^Example' -count=1 -timeout=2m .`: eight root
and one per telemetry binding have Output and actually execute. The same
selected command in both Kafka modules compiled four no-Output Examples but
reported no tests to run; none contacted a broker or entered an operator loop.
The API mode is apidiff, so prose snapshot regeneration is inapplicable; the
eight frozen baselines, SQL and every unchanged tracked input retain their bytes.
This is signature preservation, not a newly executed API qualification campaign.

`golib source-docs check` reports zero findings, `golib docs check --local` passes,
and `golib specification check` reports zero applicable modules/decisions, not
conformance certification. Immutable pinned Gitleaks scans use full redaction
and disable inline ignore annotations; baseline/full current/changed/staged
results and ordinary complete-diff review are recorded in the local handoff.
The three owned caches, helper and render output were removed immediately after
inspection, with exact absence verified. No database/broker/service, active
fuzz, race, mutation, benchmark, restart/failover or release qualification is
claimed. Publication/enforcement and exact-main CI remain coordinator-owned.

Module inventory (source/test Go files): root 37/78; gokafka 5/19;
gotelemetry 12/18; kafka 4/21; otel 11/20; outbox 2/16; queue 4/9;
competitors 0/1; postgres 5/21.

## Objective Enforcement

The checksum-pinned public Tools v2.4.0 verifier now supplies
`make source-docs-check` as a prerequisite of the existing local and CI Make
routes. Ordinary shared CI selects strict source-documentation checking and
requires its result. Native prerequisites, operation budgets, security and
release controls are retained. Hosted exact-head evidence is recorded by the
coordination ledger; configuration alone does not establish green delivery.

Continue with [Event-sourcing adoption guide](adoption.md).
