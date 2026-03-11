-- name: SaveDrawResult :execresult
INSERT INTO draw_results (giver_participant_id, receiver_participant_id, secret_friend_id, created_at, updated_at)
VALUES (?, ?, ?, NOW(), NOW());

-- name: GetDrawResultForUser :one
SELECT dr.*,
       pg.user_id                as giver_user_id,
       pr.user_id                as receiver_user_id,
       COALESCE(up.fullname, '') as receiver_fullname,
       ur.email                  as receiver_email
FROM draw_results dr
         JOIN participants pg ON dr.giver_participant_id = pg.id
         JOIN participants pr ON dr.receiver_participant_id = pr.id
         JOIN users ur ON pr.user_id = ur.id
         LEFT JOIN user_profiles up ON up.user_id = ur.id
WHERE dr.secret_friend_id = ?
  AND pg.user_id = ?;
