# DevOps — Local Postgres

Local Postgres for `libraryAPI` via Docker Compose. Schema managed by
[golang-migrate](https://github.com/golang-migrate/migrate).

## Files

- `docker-compose.yml` — Postgres 16 service, persistent volume, healthcheck.
- `migrations/` — versioned SQL (`NNNNNN_name.up.sql` / `.down.sql`).
- `embed.go` — embeds `migrations/` into the binary.

## Prerequisites

- Docker + Docker Compose

## Start

```bash
cd devops
docker compose up -d
docker compose ps        # wait for healthy
```

The DB starts empty. Migrations run automatically when the app boots
(`db.Connect` applies all pending `.up.sql`).

## Stop

```bash
docker compose down        # keep data
docker compose down -v     # wipe data (drops volume)
```

## Config

Defaults set in `docker-compose.yml`; override via env or a `.env` file
(`cp .env.example .env`):

| Var                 | Default     |
| ------------------- | ----------- |
| `POSTGRES_USER`     | `library`   |
| `POSTGRES_PASSWORD` | `library`   |
| `POSTGRES_DB`       | `librarydb` |
| `POSTGRES_PORT`     | `5432`      |

## Connect the API

`db/db.go` reads `DATABASE_URL`, falling back to the local docker DB:

```
postgresql://library:library@localhost:5432/librarydb?sslmode=disable
```

## Migrations

Applied automatically on app startup (embedded). No CLI needed to run the app.

### Add a migration

Create the next pair in `migrations/`:

```
000002_add_something.up.sql
000002_add_something.down.sql
```

Rebuild/restart the app — new `.up.sql` applied.

### Manual control (optional, needs migrate CLI)

```bash
migrate -path migrations -database "$DATABASE_URL" up
migrate -path migrations -database "$DATABASE_URL" down 1
migrate -path migrations -database "$DATABASE_URL" version
```

Install: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

> Migration state tracked in the `schema_migrations` table. If a migration
> fails mid-way it may be marked `dirty` — fix the SQL, then
> `migrate ... force <version>` before re-running.

## Inspect

```bash
docker exec -it libraryapi-postgres psql -U library -d librarydb
```
