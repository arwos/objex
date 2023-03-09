CREATE TABLE `user_group`
(
    `id`         int(10) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`    int(10) unsigned NOT NULL,
    `group_id`   int(10) unsigned NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`),
    KEY          `user_id` (`user_id`),
    KEY          `group_id` (`group_id`),
    CONSTRAINT `user_group_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `user_group_ibfk_2` FOREIGN KEY (`group_id`) REFERENCES `group` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;