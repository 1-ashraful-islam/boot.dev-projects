-- name: CreateFeed :one
INSERT INTO feeds (
  id, created_at, updated_at, user_id, url, title, description
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

