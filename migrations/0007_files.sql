CREATE TABLE `files`
(
    `id`         int(10) unsigned NOT NULL AUTO_INCREMENT,
    `storage_id` int(10) unsigned NOT NULL,
    `name`       text COLLATE utf8mb4_unicode_ci        NOT NULL,
    `hash`       varchar(40) COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` datetime                               NOT NULL,
    `updated_at` datetime                               NOT NULL,
    PRIMARY KEY (`id`),
    KEY          `storage_id` (`storage_id`),
    CONSTRAINT `files_ibfk_1` FOREIGN KEY (`storage_id`) REFERENCES `storage` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;