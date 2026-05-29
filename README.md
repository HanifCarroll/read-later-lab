# Read Later Lab

A weekend architecture experiment: a Go API server with an embedded SvelteKit frontend for a small AI-assisted read-it-later app.

The app lets you save links to articles and blog posts, capture why you saved them, and ask OpenAI to generate a brief triage card with a summary, estimated read time, recommended reading mode, best sections to read, and read/skip guidance.

## Features

- Save article/blog post URLs with a personal “why I saved this” note.
- Fetch article metadata and page text from the Go backend.
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

## Architecture

```text
cmd/server        Go HTTP server entrypoint
internal/app      domain model, article fetching, OpenAI analysis
internal/storage  JSON-file persistence for the experiment
internal/httpserver API routes and SvelteKit asset/proxy serving
ui/svelte         SvelteKit frontend
web/build         generated static frontend assets embedded by Go
```

Routes:

- `/app` — SvelteKit frontend.
- `/api/items` — JSON API consumed by SvelteKit.
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
```

`.env` is ignored by git.

## Development

```sh
make dev
```

Open <http://localhost:8080>.

Dev mode starts:

- Go server on `localhost:8080`
- SvelteKit/Vite dev server on `localhost:5173`
- Go proxying `/app/*` to Vite for HMR

## Build

```sh
make build
./server
```

The production binary serves the embedded SvelteKit build from `web/build`.

## Verification

Useful checks:

```sh
make test
make build
```
