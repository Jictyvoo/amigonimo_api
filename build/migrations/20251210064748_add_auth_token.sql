-- Modify "users" table
ALTER TABLE `users` ADD COLUMN `username` varchar(78) NOT NULL AFTER `email`, ADD COLUMN `password` varchar(76) NOT NULL AFTER `username`, ADD COLUMN `verified_at` timestamp NULL AFTER `password`, ADD COLUMN `remember_token` varchar(52) NULL AFTER `verified_at`, ADD COLUMN `verification_code` varchar(116) NULL AFTER `remember_token`, ADD COLUMN `recovery_code` varchar(27) NULL AFTER `verification_code`, ADD COLUMN `recovery_code_expires_at` timestamp NULL, ADD INDEX `user_recovery_code` (`recovery_code`), ADD UNIQUE INDEX `user_username` (`username`), ADD INDEX `user_verification_code` (`verification_code`), ADD UNIQUE INDEX `username` (`username`);
-- Create "auth_tokens" table
CREATE TABLE `auth_tokens` (
  `id` binary(16) NOT NULL DEFAULT (uuid_to_bin(uuid())),
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `token` varchar(52) NOT NULL,
  `refresh_token` char(36) NULL,
  `expires_at` timestamp NOT NULL,
  `user_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `authtoken_id` (`id`),
  INDEX `authtoken_refresh_token` (`refresh_token`),
  INDEX `authtoken_token` (`token`),
  INDEX `authtoken_user_id` (`user_id`),
  CONSTRAINT `auth_tokens_users_auth_token` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
