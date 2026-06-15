-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: ResetFeeds :exec
DELETE FROM feeds;

-- name: GetFeeds :many
SELECT feeds.*, users.name AS user_name FROM feeds LEFT JOIN users ON users.id = feeds.user_id;

-- name: GetFeed :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
--UPDATE feeds SET last_Fetched_at = current_timestamp WHERE id = $1;

-- name: GetNextFeedToFetch :one
--SELECT * FROM feeds ORDER BY last_Fetched_at NULLS FIRST LIMIT 1;
