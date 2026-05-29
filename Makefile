.PHONY: build compose-up compose-down dev db-up sqlc svelte-build test

svelte-build:
	cd ui/svelte && npm install && npm run build

build: svelte-build
	go build ./cmd/server

test:
	go test ./...

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

db-up:
	docker compose up -d postgres

compose-up:
	docker compose up --build

compose-down:
	docker compose down

dev: db-up
	@cd ui/svelte && npm install
	@sh -c 'cd ui/svelte && npm run dev' & \
		svelte_pid=$$!; \
		trap 'kill $$svelte_pid 2>/dev/null || true' INT TERM EXIT; \
		DEV_SVELTE_PROXY=http://localhost:5173 go run ./cmd/server
