# Read Later Lab

A weekend architecture experiment: a production-shaped Go API server with Postgres, sqlc, Kamal deployment, and an embedded SvelteKit frontend for a small AI-assisted read-it-later app.

The app lets you save links to articles and blog posts, capture why you saved them, and ask OpenAI to generate a brief triage card with a summary, estimated read time, recommended reading mode, best sections to read, and read/skip guidance.

## Features

- Save article/blog post URLs with a personal “why I saved this” note.
- Fetch article metadata and page text from the Go backend.
- Store items in Postgres using sqlc-generated queries.
- Run embedded migrations on app startup.
- Generate AI triage output with OpenAI:
  - brief summary
  - estimated read time
  - recommended mode: read fully, skim, reference, or skip
  - best sections to read
  - “You should read this if” bullets
  - “You should skip this if” bullets
- Triage saved items as read soon, skim later, reference, skipped, or archived.
- Run SvelteKit with Vite HMR in development.
- Build SvelteKit to static files and embed them into a single Go binary for production.
- Deploy with Docker + Kamal + a Postgres accessory.

## Architecture

```text
cmd/server                  Go HTTP server entrypoint
internal/app                domain model, article fetching, OpenAI analysis
internal/db                 sqlc schema, queries, generated Go code
internal/storage            Postgres store + embedded migrations
internal/httpserver         API routes and SvelteKit asset/proxy serving
internal/storage/migrations Postgres migrations embedded into the binary
ui/svelte                   SvelteKit frontend
web/build                   generated static frontend assets embedded by Go
config/deploy.yml           Kamal deployment config
compose.yml                 local/container Postgres + app stack
```

Routes:

- `/app` — SvelteKit frontend.
- `/api/items` — JSON API consumed by SvelteKit.
- `/healthz` — deployment health check.
- `/` — redirects to `/app`.

## Local setup

Create a local `.env` file and add your OpenAI key:

```sh
cp .env.example .env
```

Then edit `.env`:

```env
OPENAI_API_KEY=sk-your-key-here
OPENAI_MODEL=gpt-5.4-mini
ADDR=:8080
DATABASE_URL=postgres://read_later:password@127.0.0.1:55432/read_later?sslmode=disable
POSTGRES_PASSWORD=password
```

`.env` is ignored by git.

## Development

```sh
make dev
```

Open <http://localhost:8080>.

Dev mode starts:

- Postgres via Docker Compose on `localhost:5432`
- Go server on `localhost:8080`
- SvelteKit/Vite dev server on `localhost:5173`
- Go proxying `/app/*` to Vite for HMR

## Container stack

```sh
make compose-up
```

This builds the production container and starts Postgres locally.

## Build

```sh
make build
./server
```

The production binary serves the embedded SvelteKit build from `web/build` and runs DB migrations on startup.

## sqlc

After changing `internal/db/schema.sql` or `internal/db/queries/*.sql`, regenerate query code:

```sh
make sqlc
```

## Kamal deployment

The Kamal config is in `config/deploy.yml`. It deploys:

- one Go web container
- one Postgres accessory
- GHCR image: `ghcr.io/hanifcarroll/read-later-lab`

Required secrets:

```text
KAMAL_REGISTRY_PASSWORD
POSTGRES_PASSWORD
DATABASE_URL
OPENAI_API_KEY
```

For production, `DATABASE_URL` should point at the Kamal accessory, for example:

```text
postgres://read_later:<password>@read-later-lab-postgres:5432/read_later?sslmode=disable
```

## RAM checks on the VPS

```sh
kamal app exec 'ps -o pid,rss,comm -C server'
kamal app exec 'cat /sys/fs/cgroup/memory.current'
kamal accessory exec postgres 'ps aux'
ssh root@178.156.180.60 'docker stats --no-stream'
```

## Verification

Useful checks:

```sh
make test
make build
```
