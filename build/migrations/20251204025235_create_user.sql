-- Create "users" table
CREATE TABLE `users` (
  `id` binary(16) NOT NULL DEFAULT (uuid_to_bin(uuid())),
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `fullname` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `user_email` (`email`),
  UNIQUE INDEX `user_id` (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
