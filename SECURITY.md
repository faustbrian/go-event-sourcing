# Security Policy

## Reporting

Report suspected vulnerabilities privately through the GitHub security
advisory [reporting form](https://github.com/faustbrian/go-event-sourcing/security/advisories/new).
Do not open a public issue containing exploit
details, credentials, private fixtures, or affected deployment information.

Include the affected module and version, impact, reproduction, preconditions,
and any suggested mitigation. Reports are acknowledged as soon as practical;
timelines depend on severity and verification. The repository maintainer follows
the shared [vulnerability-management policy](https://github.com/faustbrian/go-library-tools/blob/main/docs/ecosystem/security/vulnerability-management.md)
for severity, response targets, private evidence, embargoes, advisories, and
coordinated publication of affected modules.

## Supported Versions

Security fixes are developed on `main` and published for the latest stable
major release line of each independently versioned module. Core, PostgreSQL,
and adapters can have different current major versions; a root-module release
does not update an adapter's minimum dependency version. Earlier major releases
remain available, but security backports are not guaranteed. Follow the affected
module's upgrade guidance and [`COMPATIBILITY.md`](COMPATIBILITY.md).

## Security Gates

Releases require isolated tests, race and hostile-input checks, exact coverage
and mutation results, `govulncheck`, secret scanning, license verification,
SBOM generation, provenance validation, and clean-consumer resolution. A
missing scanner or unavailable service is a failed gate, not a warning.

Security fixes MUST include a regression test that does not publish weaponized
details or real secrets. Credentials MUST be redacted from logs and evidence.

## Repository Assurance

The repository [safety and concurrency policy](AGENTS.md#safety-and-concurrency)
and [supply-chain policy](AGENTS.md#dependencies-and-supply-chain) define shared
trust boundaries and release requirements. Package-specific security guidance
refines those rules for its owned boundary.

The versioned [threat model](docs/threat-model.md) identifies first-party
controls, residual-risk owners, mitigations, and review triggers.
