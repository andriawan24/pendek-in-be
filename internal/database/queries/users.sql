-- name: GetUserByEmail :one
SELECT *
FROM users 
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: InsertUser :one
INSERT INTO users(
    name,
    email,
    password_hash
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET 
name = $1, email = $2, password_hash = $3, is_verified = $4
WHERE id = $5 AND deleted_at IS NULL
RETURNING *;