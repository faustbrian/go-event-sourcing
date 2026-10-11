# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected module's current docs.
2. Run `make inventory` and the narrow baseline gate for the module.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update
the affected stable decision entries, and complete the Specification Decisions
section of the pull request template. An unresolved interpretation or stale
source pin is release-blocking; peer behavior cannot silently select policy.

Required release mutation gates must finish with zero surviving viable mutants.

Do not add package-local workflows, permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use Go 1.27.2 for verification; public module minima remain Go 1.27.0.
An isolated checkout needs no sibling repositories or coordination goal files.
Begin with deterministic root tests and one active target when its decoding
boundary changes, using a disposable cache for each bounded session:

```bash
(
  set -eu
  GOCACHE="$(mktemp -d)"
  export GOCACHE
  trap 'rm -rf "${GOCACHE}"' EXIT
  GOWORK=off go test ./...
  GOWORK=off go test . -run '^$' -fuzz '^FuzzJSONCodecDecode$' -fuzztime=2s -timeout=2m
)
```

The deterministic command executes fuzz seeds/corpus, not active fuzzing.
It excludes intentional nested modules and the `integration` tag; it is not
durable-store or broker qualification. Install the public v2.1.3 binary archive
for your platform from
[go-library-tools releases](https://github.com/faustbrian/go-library-tools/releases/tag/v2.1.3),
checking `checksums.txt` against `.golib.yaml` and the archive against that
checksum set. Put the extracted `golib` binary on `PATH`, or select it through
`GOLIB`; do not build Tools from source. CI verifies this same public release
while checking out its pinned reusable workflow source.

```bash
make inventory
make check
make ci
make specification-check
golib mutation --module . --package ./memory
```

Run adjacent commands in one disposable-cache session as shown above.
`make check` and `make ci` select the bounded contract. The root memory race
runs once with a two-minute timeout: concurrent expected-version appenders and
checkpoint writers have one winner, and cancellation before write ownership
prevents a write. These assertions are not live PostgreSQL/isolation proof.
Shared safety checks retain unsafe/cgo/go:linkname refusal. Base-free checks
run seeds/corpus; changed-source CI selects causal active fuzz smoke. Optional
mutation selects one production package, one worker and a twenty-minute limit.
The all-module route preserves integration tags, whose owned fixtures include
PostgreSQL, Kafka and Valkey; it is not Docker- or external-service-free.
For a focused nested module use, for example,
`golib check --local --module adapters/queue` from this repository root.
Native nested-module checks run from that module with `GOWORK=off`; public
dependencies remain authoritative without sibling-source replacements.

Select full qualification and release rehearsal explicitly:

```bash
make check-full
make specification-check-online
golib api check
golib release dry-run --module .
```

Full checks retain all 32 declared fuzz targets at 10000x/2m each, six owned
documentation operations and eight 100ms/20m benchmark operations, mutation
imports, API baselines and every original module gate. Independently released
modules select their own module path for rehearsal. The scheduled carrier uses
the shared bounded route, not an automatic full or backend-matrix campaign.
Offline specification checking is ordinary; online authority/errata checks
are explicit. Report unavailable or failing commands; bounded checks alone
do not establish release qualification.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
