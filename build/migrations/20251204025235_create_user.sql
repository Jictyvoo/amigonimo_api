-- Create "users" table
CREATE TABLE `users` (
  `id` binary(16) NOT NULL DEFAULT (unhex(replace(uuid_v7(), '-', ''))),
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `fullname` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `username` varchar(78) NOT NULL,
  `password` varchar(76) NOT NULL,
  `verified_at` timestamp NULL,
  `remember_token` varchar(52) NULL,
  `verification_code` varchar(116) NULL,
  `recovery_code` varchar(27) NULL,
  `recovery_code_expires_at` timestamp NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `username` (`username`),
  UNIQUE INDEX `user_email` (`email`),
  UNIQUE INDEX `user_id` (`id`),
  INDEX `user_recovery_code` (`recovery_code`),
  UNIQUE INDEX `user_username` (`username`),
  INDEX `user_verification_code` (`verification_code`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
