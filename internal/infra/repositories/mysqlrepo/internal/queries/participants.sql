-- name: AddParticipant :execresult
INSERT INTO participants (id, secret_friend_id, user_id, joined_at, created_at, updated_at)
VALUES (?, ?, ?, NOW(), NOW(), NOW());

-- name: ListParticipantsBySecretFriend :many
SELECT sqlc.embed(p), u.fullname, u.email, u.username, u.id AS user_id
FROM participants p
         JOIN users u ON p.user_id = u.id
WHERE p.secret_friend_id = ?;

-- name: GetParticipantBySFAndUser :one
SELECT *
FROM participants
WHERE secret_friend_id = ?
  AND user_id = ?;

-- name: GetParticipantByID :one
SELECT *
FROM participants
WHERE id = ?;

-- name: DeleteParticipantBySFAndUser :exec
DELETE
FROM participants
WHERE secret_friend_id = ?
  AND user_id = ?;
