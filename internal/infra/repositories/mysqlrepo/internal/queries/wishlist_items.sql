-- name: AddWishlistItem :execresult
INSERT INTO wishlist_items (id, participant_id, label, comments, created_at, updated_at)
VALUES (?,
        COALESCE(
                sqlc.narg(participant_id),
                (SELECT id
                 FROM participants
                 WHERE user_id = sqlc.arg(user_id)
                   AND secret_friend_id = sqlc.arg(secret_friend_id)
                 LIMIT 1)
        ),
        ?, ?, NOW(), NOW());

-- name: RemoveWishlistItem :exec
DELETE
FROM wishlist_items
WHERE wishlist_items.id = ?
  AND wishlist_items.participant_id = COALESCE(
        sqlc.arg(participant_id),
        (SELECT id
         FROM participants
         WHERE user_id = sqlc.arg(user_id)
           AND secret_friend_id = sqlc.arg(secret_friend_id)
         LIMIT 1)
                                      );

-- name: GetWishlistByParticipant :many
SELECT *
FROM wishlist_items
WHERE wishlist_items.participant_id = COALESCE(
        sqlc.arg(participant_id),
        (SELECT id
         FROM participants
         WHERE user_id = sqlc.arg(user_id)
           AND secret_friend_id = sqlc.arg(secret_friend_id)
         LIMIT 1)
                                      )
ORDER BY created_at DESC;
