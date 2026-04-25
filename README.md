# calendar-sync

One-way sync from Microsoft 365 calendar to Google Calendar.

## Why this exists

Not all services support Microsoft 365 calendar integration. For example, [Granola](https://granola.ai) on a personal account doesn't allow connecting a work M365 calendar. This service mirrors work M365 events into a Google Calendar so that any tool supporting Google Calendar can see work events — without needing direct M365 access.

## How it works

```
┌─────────────────┐    REST API    ┌──────────────────────┐
│  calendar-sync  │ ──────────────► │    calendar-mcp       │
│  (this service) │                 │  (integration layer) │
└─────────────────┘                 └──────────────────────┘
                                              │
                                 ┌────────────┴────────────┐
                                 ▼                         ▼
                           ┌──────────┐             ┌──────────┐
                           │  M365    │             │  Google  │
                           │ Calendar │             │ Calendar │
                           └──────────┘             └──────────┘
```

1. Polls M365 for events since `last_sync`
2. For each event: create / update / skip based on content hash (title, start, end, description, location)
3. Events missing from M365 but present in mapping → deleted from Google
4. Runs every 10 minutes with overlap guard (skips tick if previous sync still running)

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `REST_API_URL` | `http://calendar-mcp:8080` | Base URL of calendar-mcp REST API |
| `API_KEY` | — | API key for calendar-mcp (required) |
| `SYNC_SOURCE` | — | Source calendar ID, e.g. `microsoft:<id>` (required) |
| `SYNC_TARGET` | — | Target calendar ID, e.g. `google:<id>` (required) |
| `STATE_FILE` | `/data/sync_state.json` | Path to persistent sync state |

To find calendar IDs, call `GET /api/calendars` on the calendar-mcp server.

## Running locally

```bash
REST_API_URL=http://localhost:8080 \
API_KEY=<key> \
SYNC_SOURCE=microsoft:<m365-calendar-id> \
SYNC_TARGET=google:<google-calendar-id> \
go run ./cmd/sync
```

## Deploy

```bash
cd personal_ai_stack/deploy/calendar-sync
docker compose up -d
```

Requires calendar-mcp on the same `infra` Docker network.
Image: `ghcr.io/dzarlax-ai/calendar-sync:latest` — built automatically on push to `main`.

## State file

```json
{
  "last_sync": "2026-04-25T10:00:00Z",
  "mappings": {
    "<m365-event-id>": { "google_id": "<google-event-id>", "hash": "sha256..." }
  }
}
```
