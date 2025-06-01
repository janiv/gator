-- name: CreatePost :one
INSERT INTO posts (created_at, updated_at, title, post_url, post_description,
published_at, feed_id) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT post_id, posts.created_at, posts.updated_at, title, post_url,
post_description, published_at, posts.feed_id FROM posts
INNER JOIN feed_follows f ON (f.feed_id = posts.feed_id)
WHERE f.user_id = $1
ORDER BY published_at
LIMIT $2;