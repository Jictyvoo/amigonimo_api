-- name: AddDenyListEntry :execresult
INSERT INTO denylists (id, participant_id, denied_user_id, created_at, updated_at)
VALUES (?, (SELECT id FROM participants WHERE user_id = ? AND secret_friend_id = ?), ?, NOW(), NOW());

-- name: AddDenyListEntryByID :execresult
INSERT INTO denylists (id, participant_id, denied_user_id, created_at, updated_at)
VALUES (?, ?, ?, NOW(), NOW());

-- name: RemoveDenyListEntry :exec
DELETE
FROM denylists
WHERE participant_id = (SELECT id FROM participants WHERE user_id = ? AND secret_friend_id = ?)
  AND denied_user_id = ?;

-- name: RemoveDenyListEntryByID :exec
DELETE
FROM denylists
WHERE participant_id = ?
  AND denied_user_id = ?;

-- name: GetDenyListByParticipant :many
SELECT d.*, u.fullname, u.email, u.username
FROM denylists d
         JOIN users u ON d.denied_user_id = u.id
WHERE d.participant_id = (SELECT id FROM participants WHERE user_id = ? AND secret_friend_id = ?);

-- name: GetDenyListByParticipantID :many
SELECT d.*, u.fullname, u.email, u.username
FROM denylists d
         JOIN users u ON d.denied_user_id = u.id
WHERE d.participant_id = ?;
