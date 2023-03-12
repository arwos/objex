CREATE TABLE `props`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT,
    `files_id`   int unsigned NOT NULL,
    `name`       varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NOT NULL,
    `value`      varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` datetime                                                      NOT NULL,
    `updated_at` datetime                                                      NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `files_id_name` (`files_id`,`name`),
    KEY          `files_id` (`files_id`),
    CONSTRAINT `props_ibfk_2` FOREIGN KEY (`files_id`) REFERENCES `files` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;