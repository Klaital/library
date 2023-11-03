CREATE TABLE `items` (
    `id` int NOT NULL,
    `code` varchar(63) NOT NULL,
    `code_type` varchar(31) DEFAULT 'UPC',
    `location` varchar(128) DEFAULT NULL,
    `title` varchar(127) DEFAULT NULL,
    `data_source` varchar(63) DEFAULT 'manual',
    `comments` varchar(255) DEFAULT NULL,
    `created` datetime DEFAULT NULL,
    `title_translated` varchar(127) DEFAULT NULL,
    PRIMARY KEY (`id`)
);
