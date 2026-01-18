-- Create "auth_tokens" table
CREATE TABLE `auth_tokens` (
  `id` binary(16) NOT NULL DEFAULT (uuid_v7()),
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `token` varchar(52) NOT NULL,
  `refresh_token` binary(16) NULL,
  `expires_at` timestamp NULL,
  `user_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `authtoken_id` (`id`),
  INDEX `authtoken_refresh_token` (`refresh_token`),
  INDEX `authtoken_token` (`token`),
  INDEX `authtoken_user_id` (`user_id`),
  CONSTRAINT `auth_tokens_users_auth_token` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
