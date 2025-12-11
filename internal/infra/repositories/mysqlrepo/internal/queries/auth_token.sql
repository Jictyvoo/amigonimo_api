-- name: GetUserByAuthToken :one
SELECT u.*
FROM users u
         INNER JOIN auth_tokens at ON u.id = at.user_id
WHERE at.token = ?;

-- name: GetAuthenticationToken :one
SELECT *
FROM auth_tokens
WHERE user_id = ?;

-- name: CheckAuthenticationByRefreshToken :one
SELECT *
FROM auth_tokens
WHERE refresh_token = ?;

-- name: UpsertAuthToken :execresult
INSERT INTO auth_tokens (id, user_id, token, refresh_token, expires_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, NOW(), NOW())
ON DUPLICATE KEY UPDATE token         = VALUES(token),
                        refresh_token = VALUES(refresh_token),
                        expires_at    = VALUES(expires_at),
                        updated_at    = NOW();
