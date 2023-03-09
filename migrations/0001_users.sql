CREATE TABLE `users`
(
    `id`         int(10) unsigned NOT NULL AUTO_INCREMENT,
    `login`      varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `passwd`     varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `acl`        tinytext COLLATE utf8mb4_unicode_ci NOT NULL,
    `lock`       tinyint(1) unsigned DEFAULT '0',
    `created_at` datetime                            NOT NULL,
    `updated_at` datetime                            NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `login` (`login`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;