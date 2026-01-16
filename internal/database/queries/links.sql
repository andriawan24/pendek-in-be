-- name: InsertLink :one
INSERT INTO links(
    original_url,
    short_code,
    custom_short_code,
    user_id,
    expired_at
) VALUES (
    $1, 
    $2, 
    $3, 
    $4, 
    $5
) 
RETURNING *;

-- name: GetRedirectLink :one
SELECT original_url FROM links WHERE (short_code = $1 OR custom_short_code = $1) AND deleted_at IS NULL;

-- name: GetLink :one
SELECT l.*, COUNT(cl.id) as counts FROM links l
LEFT JOIN click_logs cl ON cl.code = l.short_code OR cl.code = l.custom_short_code
WHERE l.user_id = $1 AND deleted_at IS NULL AND l.id = $2 
GROUP BY l.id
LIMIT 1;

-- name: GetLinks :many
SELECT l.*, COUNT(cl.id) as counts 
FROM links l
LEFT JOIN click_logs cl ON cl.code = l.short_code OR cl.code = l.custom_short_code
WHERE l.user_id = $1 AND l.deleted_at IS NULL
GROUP BY l.id
ORDER BY
  CASE WHEN @order_by::text = 'created_at' THEN l.created_at END DESC,
  CASE WHEN @order_by::text = 'updated_at' THEN l.updated_at END DESC,
  CASE WHEN @order_by::text = 'expired_at' THEN l.expired_at END DESC,
  CASE WHEN @order_by::text = 'counts' THEN COUNT(cl.id) END DESC
LIMIT $3
OFFSET $2;

-- name: UpdateLink :exec
UPDATE links SET custom_short_code = $1, original_url = $2, expired_at = $3 WHERE id = $4 AND deleted_at IS NULL;

-- name: GetTotalActiveLinks :one
SELECT COUNT(*) as total FROM links l WHERE l.user_id = $1 AND l.deleted_at IS NULL;

-- name: DeleteLink :exec
UPDATE links SET deleted_at = NOW() WHERE id = $1 AND user_id = $2;