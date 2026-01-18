-- Create "secret_friends" table
CREATE TABLE `secret_friends` (
  `id` binary(16) NOT NULL DEFAULT (uuid_v7()),
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `name` varchar(255) NOT NULL,
  `datetime` timestamp NULL,
  `location` varchar(255) NULL,
  `max_deny_list_size` tinyint unsigned NOT NULL DEFAULT 0,
  `invite_code` varchar(255) NOT NULL,
  `invite_link` varchar(255) NULL,
  `status` varchar(255) NOT NULL,
  `owner_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `secretfriend_id` (`id`),
  UNIQUE INDEX `secretfriend_invite_code` (`invite_code`),
  INDEX `secretfriend_owner_id` (`owner_id`),
  CONSTRAINT `secret_friends_users_secret_friends` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "participants" table
CREATE TABLE `participants` (
  `id` binary(16) NOT NULL DEFAULT (uuid_v7()),
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `joined_at` timestamp NULL,
  `secret_friend_id` binary(16) NOT NULL,
  `user_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `participant_id` (`id`),
  INDEX `participant_secret_friend_id` (`secret_friend_id`),
  INDEX `participant_user_id` (`user_id`),
  UNIQUE INDEX `participant_user_id_secret_friend_id` (`user_id`, `secret_friend_id`),
  CONSTRAINT `participants_secret_friends_participants` FOREIGN KEY (`secret_friend_id`) REFERENCES `secret_friends` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE,
  CONSTRAINT `participants_users_participants` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "denylists" table
CREATE TABLE `denylists` (
  `id` binary(16) NOT NULL DEFAULT (uuid_v7()),
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `participant_id` binary(16) NOT NULL,
  `denied_user_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `denylist_denied_user_id` (`denied_user_id`),
  UNIQUE INDEX `denylist_id` (`id`),
  INDEX `denylist_participant_id` (`participant_id`),
  UNIQUE INDEX `denylist_participant_id_denied_user_id` (`participant_id`, `denied_user_id`),
  CONSTRAINT `denylists_participants_denylist` FOREIGN KEY (`participant_id`) REFERENCES `participants` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE,
  CONSTRAINT `denylists_users_denied_entries` FOREIGN KEY (`denied_user_id`) REFERENCES `users` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "draw_results" table
CREATE TABLE `draw_results` (
  `id` binary(16) NOT NULL DEFAULT (uuid_v7()),
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `giver_participant_id` binary(16) NOT NULL,
  `receiver_participant_id` binary(16) NOT NULL,
  `secret_friend_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `drawresult_giver_participant_id` (`giver_participant_id`),
  UNIQUE INDEX `drawresult_id` (`id`),
  INDEX `drawresult_receiver_participant_id` (`receiver_participant_id`),
  INDEX `drawresult_secret_friend_id` (`secret_friend_id`),
  UNIQUE INDEX `drawresult_secret_friend_id_giver_participant_id` (`secret_friend_id`, `giver_participant_id`),
  UNIQUE INDEX `drawresult_secret_friend_id_receiver_participant_id` (`secret_friend_id`, `receiver_participant_id`),
  CONSTRAINT `draw_results_participants_given_results` FOREIGN KEY (`giver_participant_id`) REFERENCES `participants` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE,
  CONSTRAINT `draw_results_participants_received_results` FOREIGN KEY (`receiver_participant_id`) REFERENCES `participants` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE,
  CONSTRAINT `draw_results_secret_friends_draw_results` FOREIGN KEY (`secret_friend_id`) REFERENCES `secret_friends` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "wishlist_items" table
CREATE TABLE `wishlist_items` (
  `id` binary(16) NOT NULL DEFAULT (uuid_v7()),
  `created_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `updated_at` timestamp NOT NULL DEFAULT (current_timestamp()),
  `label` varchar(255) NOT NULL,
  `comments` longtext NULL,
  `participant_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `wishlistitem_id` (`id`),
  INDEX `wishlistitem_participant_id` (`participant_id`),
  CONSTRAINT `wishlist_items_participants_wishlist_items` FOREIGN KEY (`participant_id`) REFERENCES `participants` (`id`) ON UPDATE RESTRICT ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
