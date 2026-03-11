-- name: AddDenyListEntry :execresult
INSERT INTO denylists (id, participant_id, denied_user_id, created_at, updated_at)
VALUES (?,
        COALESCE(
                sqlc.narg(participant_id),
                (SELECT p.id
                 FROM participants p
                 WHERE p.user_id = sqlc.arg(user_id)
                   AND p.secret_friend_id = sqlc.arg(secret_friend_id)
                 LIMIT 1)
        ),
        ?, NOW(), NOW());

-- name: RemoveDenyListEntry :exec
DELETE
FROM denylists
WHERE participant_id = COALESCE(
        sqlc.arg(participant_id),
        (SELECT p.id
         FROM participants p
         WHERE p.user_id = sqlc.arg(user_id)
           AND p.secret_friend_id = sqlc.arg(secret_friend_id)
         LIMIT 1)
                       )
  AND denied_user_id = ?;

-- name: GetDenyListByParticipant :many
SELECT sqlc.embed(d), COALESCE(up.fullname, '') AS fullname, u.email, u.username, u.id AS user_id
FROM denylists d
         JOIN users u ON d.denied_user_id = u.id
         LEFT JOIN user_profiles up ON up.user_id = u.id
WHERE d.participant_id = COALESCE(
        sqlc.arg(participant_id),
        (SELECT p.id
         FROM participants p
         WHERE p.user_id = sqlc.arg(user_id)
           AND p.secret_friend_id = sqlc.arg(secret_friend_id)
         LIMIT 1)
                         );
