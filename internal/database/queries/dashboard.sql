-- name: GetTotalActiveUsers :one
SELECT COUNT(*) FROM users 
WHERE is_active = TRUE AND deleted_at IS NULL;

-- name: GetTotalLinksCreated :one
SELECT COUNT(*) FROM links 
WHERE deleted_at IS NULL;

-- name: GetGlobalTotalClicks :one
SELECT COUNT(*) FROM click_logs;
