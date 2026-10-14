# 历史实现与验收记录

以下为 2026-10-14 整理中文首页前的原始英文记录，保留具体参数、证据编号及历史状态。相对路径以仓库根目录为基准；当前状态以根目录 README 和集成计划为准。

---

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

The experimental local rc.1/rc.2 proxy archives were discarded: Windows line
conversion and archive directory entries produced noncanonical module sums.
Use `go mod download` directly against the fixed Git tag, then `go mod vendor`;
do not hand-build module ZIPs. rc.2's canonical module sum is
`h1:SunD5Oyh1le4sv3CMJTpKCa0SsKBQDqXoBoAC0aji98=` at commit
`f45ca168dee81f864533a8f77f8c203960caad1a`. Original tags/history are retained.

The SDK is not a deployment unit and contains no service-specific defaults.
`RUNTIME_URL` enables optional instance heartbeats; omitting it disables them.
Secrets and endpoint configuration are supplied by the service deployment.
