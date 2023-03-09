CREATE TABLE `providers`
(
    `id`         int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`       varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
    `code`       varchar(10) COLLATE utf8mb4_unicode_ci  NOT NULL,
    `created_at` datetime                                NOT NULL,
    `updated_at` datetime                                NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;