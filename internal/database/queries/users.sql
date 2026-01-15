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
name = $1, email = $2, password_hash = $3, is_verified = $4, profile_image_url = $5
WHERE id = $6 AND deleted_at IS NULL
RETURNING *;

-- name: GetUserByGoogleID :one
SELECT *
FROM users
WHERE google_id = $1 AND deleted_at IS NULL;

-- name: InsertUserWithGoogle :one
INSERT INTO users(
    name,
    email,
    google_id,
    is_verified,
    profile_image_url
) VALUES (
    $1,
    $2,
    $3,
    TRUE,
    $4
)
RETURNING *;