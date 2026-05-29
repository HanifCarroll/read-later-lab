-- name: ListItems :many
SELECT id, url, title, site, author, published_at, saved_reason, summary, read_time_minutes, recommended_mode, best_sections, read_if, skip_if, status, created_at, updated_at
FROM items
ORDER BY created_at DESC;

-- name: CreateItem :one
INSERT INTO items (
  id, url, title, site, author, published_at, saved_reason, summary, read_time_minutes, recommended_mode, best_sections, read_if, skip_if, status, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
)
RETURNING id, url, title, site, author, published_at, saved_reason, summary, read_time_minutes, recommended_mode, best_sections, read_if, skip_if, status, created_at, updated_at;

-- name: UpdateItemStatus :one
UPDATE items
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING id, url, title, site, author, published_at, saved_reason, summary, read_time_minutes, recommended_mode, best_sections, read_if, skip_if, status, created_at, updated_at;
