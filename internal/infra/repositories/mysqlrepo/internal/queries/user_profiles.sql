-- name: CreateUserProfile :execresult
INSERT INTO user_profiles (id, user_id, fullname, nickname, image_link, birthday, address, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW());

-- name: GetUserProfileByUserID :one
SELECT *
FROM user_profiles
WHERE user_id = ?;

-- name: UpsertUserProfile :execresult
INSERT INTO user_profiles (id, user_id, fullname, nickname, image_link, birthday, address, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
ON DUPLICATE KEY UPDATE fullname   = VALUES(fullname),
                        nickname   = VALUES(nickname),
                        image_link = VALUES(image_link),
                        birthday   = VALUES(birthday),
                        address    = VALUES(address),
                        updated_at = NOW();
