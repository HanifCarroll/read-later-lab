package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"go-embedded-js/internal/app"
	"go-embedded-js/internal/httpserver"
	"go-embedded-js/internal/storage"
	"go-embedded-js/web"
)

func main() {
	loadDotEnv()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, envOr("DATABASE_URL", "postgres://read_later:password@127.0.0.1:55432/read_later?sslmode=disable"))
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer pool.Close()
	if err := storage.Migrate(ctx, pool); err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	service := app.NewService(storage.NewPostgresStore(pool), app.OpenAIAnalyzer{
		APIKey: os.Getenv("OPENAI_API_KEY"),
		Model:  envOr("OPENAI_MODEL", "gpt-5.4-mini"),
	})

	handler := httpserver.New(httpserver.Config{
		Service:        service,
		Assets:         web.Assets,
		DevSvelteProxy: os.Getenv("DEV_SVELTE_PROXY"),
	})

	addr := envOr("ADDR", ":8080")
	log.Printf("read-later listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loadDotEnv() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("warning: could not load .env: %v", err)
	}
}
