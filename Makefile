.PHONY: dev build test svelte-build

svelte-build:
	cd ui/svelte && npm install && npm run build

build: svelte-build
	go build ./cmd/server

test:
	go test ./...

dev:
	@cd ui/svelte && npm install
	@sh -c 'cd ui/svelte && npm run dev' & \
		svelte_pid=$$!; \
		trap 'kill $$svelte_pid 2>/dev/null || true' INT TERM EXIT; \
		DEV_SVELTE_PROXY=http://localhost:5173 go run ./cmd/server
