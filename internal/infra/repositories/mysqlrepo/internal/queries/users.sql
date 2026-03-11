-- name: CreateUser :execresult
INSERT INTO users (id, email, username, password, verification_code, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, NOW(), NOW());

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = ?;

-- name: GetUserByVerificationCode :one
SELECT *
FROM users
WHERE verification_code = ?;

-- name: GetUserByRecovery :one
SELECT *
FROM users
WHERE email = ?
  AND recovery_code = ?
  AND recovery_code_expires_at >= ?;

-- name: GetUserByEmailOrUsername :one
SELECT *
FROM users
WHERE email = ?
   OR username = ?
LIMIT 1;

-- name: SetUserVerified :exec
UPDATE users
SET verified_at = NOW(),
    updated_at  = NOW()
WHERE id = ?;

-- name: SetRecoveryCode :exec
UPDATE users
SET recovery_code            = ?,
    recovery_code_expires_at = ?,
    updated_at               = NOW()
WHERE id = ?;

-- name: UpdatePassword :exec
UPDATE users
SET password   = ?,
    updated_at = NOW()
WHERE id = ?;
