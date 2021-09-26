use rabbit_gather;

show tables;

drop table if exists article_tag;
drop table if exists articles;
drop table if exists article_tags;
drop table if exists permissions;

drop table if exists tag_types;
drop table if exists tags;
drop table if exists permissions;
drop table if exists role_permission;
drop table if exists roles;
drop table if exists users;
drop table if exists role;


show create table rabbit_gather.users;
CREATE TABLE `users`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at`    datetime(3)         DEFAULT NULL,
    `updated_at`    datetime(3)         DEFAULT NULL,
    `deleted_at`    datetime(3)         DEFAULT NULL,
    `role_id`       bigint(20) unsigned DEFAULT NULL,
    `first_name`    varchar(50)         DEFAULT NULL,
    `middle_name`   varchar(50)         DEFAULT NULL,
    `last_name`     varchar(50)         DEFAULT NULL,
    `mobile`        varchar(15)         DEFAULT NULL,
    `email`         varchar(50)         DEFAULT NULL,
    `password_hash` char(60)            NOT NULL,
    `password_salt` char(24)            NOT NULL,
    `registered_at` datetime(3)         DEFAULT NULL,
    `last_login`    datetime(3)         DEFAULT NULL,
    `intro`         tinytext            DEFAULT NULL,
    `profile`       text                DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_users_deleted_at` (`deleted_at`),
    KEY `fk_roles_user` (`role_id`),
    CONSTRAINT `fk_roles_user` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
CREATE TABLE `roles`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at`  datetime(3)                  DEFAULT NULL,
    `updated_at`  datetime(3)                  DEFAULT NULL,
    `deleted_at`  datetime(3)                  DEFAULT NULL,
    `title`       varchar(75)         NOT NULL,
    `slug`        varchar(100)        NOT NULL,
    `description` tinytext                     DEFAULT NULL,
    `active`      tinyint(1)          NOT NULL DEFAULT 0,
    `content`     text                         DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_roles_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE `users`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at`    datetime(3)         DEFAULT NULL,
    `updated_at`    datetime(3)         DEFAULT NULL,
    `deleted_at`    datetime(3)         DEFAULT NULL,
    `first_name`    varchar(50)         DEFAULT NULL,
    `middle_name`   varchar(50)         DEFAULT NULL,
    `last_name`     varchar(50)         DEFAULT NULL,
    `mobile`        varchar(15)         DEFAULT NULL,
    `email`         varchar(50)         DEFAULT NULL,
    `password_hash` char(60)            NOT NULL,
    `registered_at` datetime(3)         DEFAULT NULL,
    `last_login`    datetime(3)         DEFAULT NULL,
    `intro`         tinytext            DEFAULT NULL,
    `profile`       text                DEFAULT NULL,
    `password_salt` char(24)            NOT NULL,
    `role_id`       bigint(20) unsigned DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;



CREATE TABLE `article_tag`
(
    `article_id` bigint(20) unsigned NOT NULL,
    `tag_id`     bigint(20) unsigned NOT NULL,
    PRIMARY KEY (`article_id`, `tag_id`),
    KEY `fk_article_tag_tag` (`tag_id`),
    CONSTRAINT `fk_article_tag_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`),
    CONSTRAINT `fk_article_tag_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE `tags`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `name`       char(24)            NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_tags_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
CREATE TABLE `articles`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `title`      varchar(75)         NOT NULL,
    `content`    mediumtext          NOT NULL,
    `coords`     point               NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_articles_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
CREATE TABLE `permissions`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at`  datetime(3)                  DEFAULT NULL,
    `updated_at`  datetime(3)                  DEFAULT NULL,
    `deleted_at`  datetime(3)                  DEFAULT NULL,
    `title`       varchar(75)         NOT NULL,
    `slug`        longtext                     DEFAULT NULL,
    `description` tinytext                     DEFAULT NULL,
    `active`      tinyint(1)          NOT NULL DEFAULT 0,
    `content`     text                         DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_permissions_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
show create table users;
CREATE TABLE `users`
(
    `id`             bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at`     datetime(3) DEFAULT NULL,
    `updated_at`     datetime(3) DEFAULT NULL,
    `deleted_at`     datetime(3) DEFAULT NULL,
    `first_name`     varchar(50) DEFAULT NULL,
    `middle_name`    varchar(50) DEFAULT NULL,
    `last_name`      varchar(50) DEFAULT NULL,
    `mobile`         varchar(15) DEFAULT NULL,
    `email`          varchar(50) DEFAULT NULL,
    `password_hash`  char(60)            NOT NULL,
    `registered_at`  datetime(3) DEFAULT NULL,
    `last_login`     datetime(3) DEFAULT NULL,
    `intro`          tinytext    DEFAULT NULL,
    `profile`        text        DEFAULT NULL,
    `password_salt5` char(24)            NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

show create table article_tags;
CREATE TABLE `articles`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `title`      varchar(75)         NOT NULL,
    `content`    mediumtext          NOT NULL,
    `coords`     point               NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_articles_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


CREATE TABLE `article_tags`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `name`       char(24)            NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_article_tags_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_articles_article_tag` FOREIGN KEY (`id`) REFERENCES `articles` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
# CREATE TABLE `role_permission`
# (
#     `permission_id` bigint(20) unsigned NOT NULL,
#     `role_id`       bigint(20) unsigned NOT NULL,
#     PRIMARY KEY (`permission_id`, `role_id`),
#     KEY `fk_role_permission_role` (`role_id`),
#     CONSTRAINT `fk_role_permission_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`),
#     CONSTRAINT `fk_role_permission_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
# ) ENGINE = InnoDB
#   DEFAULT CHARSET = utf8mb4;


CREATE TABLE `role_permissions`
(
    CONSTRAINT `fk_role_permissions_role`
        FOREIGN KEY (`) REFERENCES `roles`(`id`),CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`) REFERENCES `permissions` (`id`)
) ENGINE = InnoDB;

show create table users;

CREATE TABLE `users`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `first_name` longtext    DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE `users`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `first_name` longtext    DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;