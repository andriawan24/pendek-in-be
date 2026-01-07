-- name: InsertClickLog :one
INSERT INTO click_logs (
    code,
    ip_address,
    user_agent,
    referrer,
    country,
    traffic,
    device_type,
    browser
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetByDateRange :many
SELECT 
    DATE_TRUNC('day', cl.clicked_at)::timestamp AS date,
    COUNT(*) AS total_click
FROM click_logs cl
LEFT JOIN links l ON l.short_code = cl.code OR l.custom_short_code = cl.code
WHERE cl.clicked_at BETWEEN @from_date::timestamp AND @to_date::timestamp AND l.user_id = $1
GROUP BY DATE_TRUNC('day', cl.clicked_at)
ORDER BY date ASC;