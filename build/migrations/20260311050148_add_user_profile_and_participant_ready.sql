-- Modify "participants" table
ALTER TABLE `participants` ADD COLUMN `is_ready` bool NOT NULL DEFAULT 0 AFTER `joined_at`;
-- Modify "users" table
ALTER TABLE `users` DROP COLUMN `fullname`;
-- Create "user_profiles" table
CREATE TABLE `user_profiles` (
  `id` binary(16) NOT NULL DEFAULT (unhex(replace(uuid_v7(), '-', ''))),
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `fullname` varchar(255) NULL,
  `nickname` varchar(120) NULL,
  `image_link` varchar(2048) NULL,
  `birthday` timestamp NULL,
  `address` varchar(255) NULL,
  `user_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `userprofile_birthday` (`birthday`),
  UNIQUE INDEX `userprofile_id` (`id`),
  INDEX `userprofile_nickname` (`nickname`),
  UNIQUE INDEX `userprofile_user_id` (`user_id`),
  UNIQUE INDEX `user_id` (`user_id`),
  CONSTRAINT `user_profiles_users_profile` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
