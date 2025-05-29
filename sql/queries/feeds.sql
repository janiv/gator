-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1 LIMIT 1;

-- name: MarkFeedFetched :exec
UPDATE feeds SET last_fetched = $1, updated_at = $1 WHERE
id = $2;

-- name: GetNextFeedToFetch :one
SELECT id, url FROM feeds ORDER BY last_fetched NULLS FIRST LIMIT 1;
