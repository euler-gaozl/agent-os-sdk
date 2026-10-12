# Repository boundaries

- This repository is the shared protocol/transport module, not a business service.
- Keep SQL, table ownership, route permission policies and migrations in owners.
- Do not import `agent-os/internal`, sibling source or unpublished local replaces.
- Preserve wire compatibility; publish new fixed versions for changed behavior.
- Never include credentials, real service configuration or personal data.
- Use Go; tests must exercise boundaries/failures rather than mirror source.
- Verify with GOWORK=off, go test and go vet. Report skipped integration work.
