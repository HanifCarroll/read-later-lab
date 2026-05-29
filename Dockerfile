FROM node:24-slim AS frontend
WORKDIR /app
COPY ui/svelte/package*.json ./ui/svelte/
RUN cd ui/svelte && npm install
COPY ui/svelte ./ui/svelte
COPY web ./web
RUN cd ui/svelte && npm run build

FROM golang:1.26-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/build ./web/build
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM debian:bookworm-slim AS release
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /out/server ./server
EXPOSE 8080
CMD ["./server"]
