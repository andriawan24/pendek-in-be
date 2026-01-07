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

-- name: GetTotalClicks :one
WITH click_counts AS (
    SELECT l.id, COUNT(cl.*) AS total_clicks
    FROM links l
    LEFT JOIN click_logs cl ON l.short_code = cl.code OR l.custom_short_code = cl.code
    WHERE l.user_id = $1 AND l.deleted_at IS NULL AND
    cl.clicked_at BETWEEN @from_date::timestamp AND @to_date::timestamp AND l.user_id = $1
    GROUP BY l.id
)
SELECT COALESCE(SUM(total_clicks), 0) AS total FROM click_counts;

-- name: GetByDateRange :many
SELECT 
    DATE_TRUNC('day', cl.clicked_at)::timestamp AS date,
    COUNT(*) AS total_click
FROM click_logs cl
LEFT JOIN links l ON l.short_code = cl.code OR l.custom_short_code = cl.code
WHERE cl.clicked_at BETWEEN @from_date::timestamp AND @to_date::timestamp AND l.user_id = $1
GROUP BY DATE_TRUNC('day', cl.clicked_at)
ORDER BY date ASC;

-- name: GetDeviceBreakdown :many
SELECT 
    COALESCE(cl.device_type, 'Unknown') AS device_type,
    COUNT(*) AS total
FROM click_logs cl
LEFT JOIN links l ON l.short_code = cl.code OR l.custom_short_code = cl.code
WHERE cl.clicked_at BETWEEN @from_date::timestamp AND @to_date::timestamp AND l.user_id = $1
GROUP BY cl.device_type
ORDER BY total DESC;

-- name: GetTopCountries :many
SELECT 
    COALESCE(cl.country, 'Unknown') AS country,
    COUNT(*) AS total
FROM click_logs cl
LEFT JOIN links l ON l.short_code = cl.code OR l.custom_short_code = cl.code
WHERE cl.clicked_at BETWEEN @from_date::timestamp AND @to_date::timestamp AND l.user_id = $1
GROUP BY cl.country
ORDER BY total DESC
LIMIT 10;

-- name: GetTrafficSources :many
SELECT 
    COALESCE(cl.traffic, 'Direct') AS traffic_source,
    COUNT(*) AS total
FROM click_logs cl
LEFT JOIN links l ON l.short_code = cl.code OR l.custom_short_code = cl.code
WHERE cl.clicked_at BETWEEN @from_date::timestamp AND @to_date::timestamp AND l.user_id = $1
GROUP BY cl.traffic
ORDER BY total DESC;

-- name: GetBrowserUsage :many
SELECT 
    COALESCE(cl.browser, 'Unknown') AS browser,
    COUNT(*) AS total
FROM click_logs cl
LEFT JOIN links l ON l.short_code = cl.code OR l.custom_short_code = cl.code
WHERE cl.clicked_at BETWEEN @from_date::timestamp AND @to_date::timestamp AND l.user_id = $1
GROUP BY cl.browser
ORDER BY total DESC;