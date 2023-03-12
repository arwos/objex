CREATE TABLE `files`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT,
    `storage_id` int unsigned NOT NULL,
    `name`       text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `hash`       varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` datetime                                                     NOT NULL,
    `updated_at` datetime                                                     NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `hash` (`hash`),
    KEY          `storage_id` (`storage_id`),
    CONSTRAINT `files_ibfk_1` FOREIGN KEY (`storage_id`) REFERENCES `storage` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;