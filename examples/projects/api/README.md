# Mocked notes API

Small multi-file backend demo using the shipped SF-3 `@std` surface: SQLite
(`@std/db`), JSON (`@std/json`), env (`@std/env`), and HTTP (`@std/http`).

No external cloud services — default DB is in-memory SQLite (`API_DB`, default
`:memory:`).

## Layout

| File | Role |
| --- | --- |
| `store.vol` | SQLite open/migrate + notes CRUD |
| `handlers.vol` | Route `handle(req)` + `boot()` |
| `main.vol` | In-process self-test with mock requests (no `listen`) |
| `server.vol` | Real `listen` server (blocks until killed) |

## Endpoints

| Method | Path | Behavior |
| --- | --- | --- |
| `GET` | `/health` | `{"ok":true}` |
| `GET` | `/notes` | JSON array of notes |
| `POST` | `/notes` | Body `{"title":"…"}` → create (`201`) |
| `GET` | `/notes/:id` | One note or `404` |
| `DELETE` | `/notes/:id` | `{"deleted":true}` or `404` |

## Run

Self-test (what `examples_test` covers):

```text
go run ./cmd/vol run ./examples/projects/api/main.vol
```

HTTP server (blocks):

```text
API_ADDR=127.0.0.1:8787 API_DB=./notes.db go run ./cmd/vol run ./examples/projects/api/server.vol
```

```text
curl -s http://127.0.0.1:8787/health
curl -s -X POST http://127.0.0.1:8787/notes -d '{"title":"ship it"}'
curl -s http://127.0.0.1:8787/notes
curl -s http://127.0.0.1:8787/notes/1
curl -s -X DELETE http://127.0.0.1:8787/notes/1
```
