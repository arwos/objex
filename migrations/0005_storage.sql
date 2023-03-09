CREATE TABLE `storage`
(
    `id`          int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `lifetime`    int(10) unsigned NOT NULL DEFAULT '0',
    `provider_id` int(10) unsigned NOT NULL,
    `created_at`  datetime                                NOT NULL,
    `updated_at`  datetime                                NOT NULL,
    PRIMARY KEY (`id`),
    KEY           `provider_id` (`provider_id`),
    CONSTRAINT `storage_ibfk_1` FOREIGN KEY (`provider_id`) REFERENCES `providers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;