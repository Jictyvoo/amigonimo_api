-- name: CreateSecretFriend :execresult
INSERT INTO secret_friends (id, name, datetime, location, invite_code, status, owner_id, max_deny_list_size, created_at,
                            updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW());

-- name: GetSecretFriendByID :one
SELECT *
FROM secret_friends
WHERE id = ?;

-- name: UpdateSecretFriend :exec
UPDATE secret_friends
SET name       = ?,
    datetime   = ?,
    location   = ?,
    status     = ?,
    updated_at = NOW()
WHERE id = ?;

-- name: ListSecretFriends :many
SELECT *
FROM secret_friends
WHERE owner_id = ?
   OR id IN (SELECT secret_friend_id FROM participants p WHERE p.user_id = ?)
ORDER BY created_at DESC;

-- name: GetSecretFriendByInviteCode :one
SELECT *
FROM secret_friends
WHERE invite_code = ?;
