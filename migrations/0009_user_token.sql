CREATE TABLE `user_token`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT,
    `token`      varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `user_id`    int unsigned NOT NULL,
    `created_at` datetime                                                     NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`token`),
    KEY          `user_id` (`user_id`),
    CONSTRAINT `user_token_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;