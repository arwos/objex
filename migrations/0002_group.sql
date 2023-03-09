CREATE TABLE `group`
(
    `id`         int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`       varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` datetime                                NOT NULL,
    `updated_at` datetime                                NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;