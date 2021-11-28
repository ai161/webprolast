-- Table for tasks
DROP TABLE IF EXISTS `tasks`;

CREATE TABLE `tasks` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `title` varchar(50) NOT NULL,
    `is_done` boolean NOT NULL DEFAULT b'0',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8mb4;

CREATE TABLE `accounts`(
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `name` varchar(20) NOT NULL,
    `password` varchar(256) NOT NULL,
    `birth` varchar(20) NOT NULL DEFAULT "",
    PRIMARY KEY (`id`)
)DEFAULT CHARSET=utf8mb4;

CREATE TABLE `task_owners`(
    `user_id` bigint(20) NOT NULL,
    `task_id` bigint(20) NOT NULL,
    PRIMARY KEY (`user_id`,`task_id`)
)DEFAULT CHARSET=utf8mb4;