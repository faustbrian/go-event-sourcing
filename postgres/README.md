# event-sourcing PostgreSQL

`postgres` is the independently releasable PostgreSQL adapter for
`github.com/faustbrian/go-event-sourcing`. Installing the core module
does not add `pgx` or database dependencies.

## Install

```sh
go get github.com/faustbrian/go-event-sourcing/postgres@v1
```

## Quick start

```go
source, err := migrations.NewFSSource(eventpostgres.Migrations(), ".")
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
- [Go API reference](https://pkg.go.dev/github.com/faustbrian/go-event-sourcing/postgres)
- [Parent package documentation](../docs/README.md)

## Compatibility and support

This module follows Semantic Versioning. Report vulnerabilities through the
[parent security policy](../SECURITY.md).

## License

MIT. See [LICENSE](LICENSE).
