-- name: CreateUser :execresult
INSERT INTO users (id, fullname, email, created_at, updated_at)
VALUES (?, ?, ?, NOW(), NOW());

-- name: GetUserByID :one
SELECT id, fullname, email, created_at, updated_at
FROM users
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT id, fullname, email, created_at, updated_at
FROM users
WHERE email = ?;
