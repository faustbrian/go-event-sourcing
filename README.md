# Event Sourcing

[![CI](https://github.com/faustbrian/go-event-sourcing/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-event-sourcing/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-event-sourcing/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-event-sourcing.svg)](https://pkg.go.dev/github.com/faustbrian/go-event-sourcing)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-event-sourcing?sort=semver)](https://github.com/faustbrian/go-event-sourcing/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`event-sourcing` is a pragmatic event-sourcing library for Go under active
development. It is designed around three independently replaceable
responsibilities:

1. an aggregate repository;
2. an immutable event store; and
3. a message dispatcher.

The design is inspired by EventSauce while using explicit Go composition,
small consumer-owned interfaces, and `context.Context` at I/O boundaries. It
does not require CQRS, a command bus, a query bus, a queue, an outbox, a
framework, reflection-based handler discovery, or code generation.

The package is a pre-release candidate. Its documented core and adapter
capabilities are implemented, but applications should pin a reviewed revision
until the first versioned release and its complete release evidence are
published.

## Quickstart

See the complete [five-minute quickstart](docs/quickstart.md) for one aggregate
using the conformant in-memory store and the same repository boundary used by
durable adapters.

## Design documents

- [Five-minute quickstart](docs/quickstart.md)
- [Runnable workflow examples](docs/examples.md)
- [Adoption guide and anti-patterns](docs/adoption.md)
- [Aggregate modeling](docs/aggregates.md)
- [Aggregate identifiers and UUID encoding](docs/identifiers.md)
- [Architecture](docs/architecture.md)
- [Aggregate lifecycle](docs/lifecycle.md)
- [Event messages and metadata](docs/messages.md)
- [Serialization and schema evolution](docs/serialization.md)
- [Installation and package map](docs/installation.md)
- [Learning event sourcing](docs/learning.md)
- [Frequently asked questions](docs/faq.md)
- [Security, privacy, and compliance](docs/security.md)
- [Troubleshooting](docs/troubleshooting.md)
- [Glossary](docs/glossary.md)
- [Release notes and compatibility](docs/release-notes.md)
- [Release hardening findings](docs/release-audit.md)
- [Dispatcher and consumer semantics](docs/dispatch.md)
- [Code generation decision](docs/code-generation.md)
- [Anti-corruption translation](docs/translation.md)
- [EventSauce 3.9.1 compatibility matrix](docs/compatibility/eventsauce-3.9.1.md)
- [EventSauce-to-Go migration guide](docs/migration-eventsauce.md)
- [Public API design](docs/design/public-api.md)
- [Package and adapter boundaries](docs/design/package-boundaries.md)
- [Aggregate scenario testing](docs/testing.md)
- [Snapshot storage](docs/snapshots.md)
- [Projection and replay foundations](docs/projections.md)
- [Process managers](docs/process-managers.md)
- [Database structure and capacity](docs/database-structure.md)
- [Performance and capacity planning](docs/performance.md)
- [PostgreSQL adapter](postgres/README.md)
- [Transactional outbox integration](docs/outbox.md)
- [Custom outbox boundary](docs/custom-outbox.md)
- [Kafka integration](docs/kafka.md)
- [Compatible queue integration](docs/queue.md)
- [OpenTelemetry integration](docs/telemetry.md)

## When to use event sourcing

See the [adoption guide](docs/adoption.md) for the decision criteria, bounded
adoption path, prerequisites, anti-patterns, and exit criteria. Conventional
state persistence is explicitly recommended when event history does not
justify its modeling, evolution, replay, privacy, and operational costs.

## Status

The stable v1 API follows semantic-versioning compatibility. The
compatibility matrix distinguishes implemented, excluded, and externally owned
capabilities; release readiness still depends on the complete repository gates
described in the release notes.

## Ecosystem

Use the [Golib documentation portal](https://github.com/faustbrian/golib/blob/main/docs/index.md)
to choose companion packages, supported stacks, recipes, and operations guidance.
