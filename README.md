# Agent OS SDK

Shared versioned wire contracts, bounded HTTP transport, MCP transport/security,
catalog DTOs and configurable Go process hosting. Extracted from integration
commit `400afaea7ed6c9c4aa47b8002839efae2d6bf55d`.

There are no application SQL statements, migrations, record repositories, role
authorization tables or role-to-schema probes in this module. Each entrypoint
supplies `servicehost.Options` with its own readiness and authorization policy.
Database connection pooling is infrastructure; business queries remain in the
owning repository. `contracts` defines the existing versioned wire structures
and input invariants; it does not execute business workflows.

## Verification

Go 1.27.2; run `GOWORK=off go test -count=1 ./...` and `go vet ./...`.
CI tests Linux (with the race detector) and Windows. No database or model key is
needed for SDK tests. Service integration tests belong to the owning services.

## Private distribution

Consumers pin a module version and vendor the required dependency sources,
including their licenses. Their CI builds a single checkout with `GOWORK=off`
and `-mod=vendor`; no personal GitHub token or adjacent `replace` path is needed.
The integration release tooling records the exact SDK commit/module checksum.
Changes require a new version; never overwrite a released module version.

The SDK is not a deployment unit and contains no service-specific defaults.
`RUNTIME_URL` enables optional instance heartbeats; omitting it disables them.
Secrets and endpoint configuration are supplied by the service deployment.
