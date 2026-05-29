package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-embedded-js/internal/app"
	db "go-embedded-js/internal/db/generated"
)

type PostgresStore struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool, queries: db.New(pool)}
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

func (s *PostgresStore) List(ctx context.Context) ([]app.Item, error) {
	rows, err := s.queries.ListItems(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]app.Item, 0, len(rows))
	for _, row := range rows {
		item, err := mapItem(row)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *PostgresStore) Create(ctx context.Context, item app.Item) (app.Item, error) {
	row, err := s.queries.CreateItem(ctx, db.CreateItemParams{
		ID:              item.ID,
		Url:             item.URL,
		Title:           item.Title,
		Site:            item.Site,
		Author:          item.Author,
		PublishedAt:     item.PublishedAt,
		SavedReason:     item.SavedReason,
		Summary:         item.Summary,
		ReadTimeMinutes: int32(item.ReadTimeMinutes),
		RecommendedMode: string(item.RecommendedMode),
		BestSections:    item.BestSections,
		ReadIf:          item.ReadIf,
		SkipIf:          item.SkipIf,
		Status:          string(item.Status),
		CreatedAt:       timestamptz(item.CreatedAt),
		UpdatedAt:       timestamptz(item.UpdatedAt),
	})
	if err != nil {
		return app.Item{}, err
	}
	return mapItem(row)
}

func (s *PostgresStore) UpdateStatus(ctx context.Context, id string, status app.ItemStatus) (app.Item, error) {
	row, err := s.queries.UpdateItemStatus(ctx, db.UpdateItemStatusParams{ID: id, Status: string(status)})
	if err != nil {
		return app.Item{}, err
	}
	return mapItem(row)
}

func mapItem(row db.Item) (app.Item, error) {
	createdAt, err := timeFrom(row.CreatedAt)
	if err != nil {
		return app.Item{}, fmt.Errorf("created_at: %w", err)
	}
	updatedAt, err := timeFrom(row.UpdatedAt)
	if err != nil {
		return app.Item{}, fmt.Errorf("updated_at: %w", err)
	}
	return app.Item{
		ID:              row.ID,
		URL:             row.Url,
		Title:           row.Title,
		Site:            row.Site,
		Author:          row.Author,
		PublishedAt:     row.PublishedAt,
		SavedReason:     row.SavedReason,
		Summary:         row.Summary,
		ReadTimeMinutes: int(row.ReadTimeMinutes),
		RecommendedMode: app.ReadingMode(row.RecommendedMode),
		BestSections:    row.BestSections,
		ReadIf:          row.ReadIf,
		SkipIf:          row.SkipIf,
		Status:          app.ItemStatus(row.Status),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}

func timestamptz(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: value, Valid: true}
}

func timeFrom(value pgtype.Timestamptz) (time.Time, error) {
	if !value.Valid {
		return time.Time{}, fmt.Errorf("timestamp is null")
	}
	return value.Time, nil
}
