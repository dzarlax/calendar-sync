# calendar-sync

One-way M365 → Google Calendar sync via calendar-mcp REST API.

## Commands

```bash
go build -o sync ./cmd/sync   # Build
go test -race -count=1 ./...  # Tests
go vet ./...                  # Lint
REST_API_URL=http://localhost:8080 API_KEY=test SYNC_SOURCE=microsoft:<id> SYNC_TARGET=google:<id> ./sync  # Run locally
```

## Architecture

- `cmd/sync/` — entry point, config, signal handling
- `internal/sync/rest_client.go` — HTTP client for calendar-mcp REST API
- `internal/sync/syncer.go` — core logic: poll M365 → create/update/delete Google
- `internal/sync/state.go` — sync cursor + event ID mapping (atomic JSON writes)
- `internal/sync/scheduler.go` — 10-min ticker with overlap guard
- `internal/sync/hash.go` — SHA-256 change detection (title, start, end, description, location)

## Deploy

Config: `personal_ai_stack/deploy/calendar-sync/`
Image: `ghcr.io/dzarlax/calendar-sync:latest`
