# create database `rabbit_gather` CHARACTER SET = utf8mb4;
# CREATE USER 'rabbit_gather'@'%' IDENTIFIED BY '5678';
# GRANT ALL ON rabbit_gather.* TO 'rabbit_gather'@'%';
# FLUSH PRIVILEGES;
#
# use rabbit_gather;
# show tables;
# drop table if exists role;
# show create table role;
#
# CREATE TABLE `role`
# (
#     `id`          BIGINT       NOT NULL AUTO_INCREMENT,
#     `title`       VARCHAR(75)  NOT NULL,
#     `slug`        VARCHAR(100) NOT NULL,
#     `description` TINYTEXT     NULL,
#     `active`      TINYINT(1)   NOT NULL DEFAULT 0,
#     `createdAt`   DATETIME     NOT NULL,
#     `updatedAt`   DATETIME     NULL     DEFAULT NULL,
#     `content`     TEXT         NULL     DEFAULT NULL,
#     PRIMARY KEY (`id`),
#     UNIQUE INDEX `uq_slug` (`slug` ASC)
# );
#
# CREATE TABLE `permission`
# (
#     `id`          BIGINT       NOT NULL AUTO_INCREMENT,
#     `title`       VARCHAR(75)  NOT NULL,
#     `slug`        VARCHAR(100) NOT NULL,
#     `description` TINYTEXT     NULL,
#     `active`      TINYINT(1)   NOT NULL DEFAULT 0,
#     `createdAt`   DATETIME     NOT NULL,
#     `updatedAt`   DATETIME     NULL     DEFAULT NULL,
#     `content`     TEXT         NULL     DEFAULT NULL,
#     PRIMARY KEY (`id`),
#     UNIQUE INDEX `uq_slug` (`slug` ASC)
# );
#
# CREATE TABLE `role_permission`
# (
#     `roleId`       BIGINT   NOT NULL,
#     `permissionId` BIGINT   NOT NULL,
#     `createdAt`    DATETIME NOT NULL,
#     `updatedAt`    DATETIME NULL,
#     PRIMARY KEY (`roleId`, `permissionId`),
#     INDEX `idx_rp_role` (`roleId` ASC),
#     INDEX `idx_rp_permission` (`permissionId` ASC),
#     CONSTRAINT `fk_rp_role`
#         FOREIGN KEY (`roleId`)
#             REFERENCES `role` (`id`)
#             ON DELETE NO ACTION
#             ON UPDATE NO ACTION,
#     CONSTRAINT `fk_rp_permission`
#         FOREIGN KEY (`permissionId`)
#             REFERENCES `permission` (`id`)
#             ON DELETE NO ACTION
#             ON UPDATE NO ACTION
# );
#
# CREATE TABLE `user`
# (
#     `id`           BIGINT      NOT NULL AUTO_INCREMENT,
#     `firstName`    VARCHAR(50) NULL DEFAULT NULL,
#     `middleName`   VARCHAR(50) NULL DEFAULT NULL,
#     `lastName`     VARCHAR(50) NULL DEFAULT NULL,
#     `mobile`       VARCHAR(15) NULL,
#     `email`        VARCHAR(50) NULL,
#     `passwordHash` VARCHAR(32) NOT NULL,
#     `passwordHash` VARCHAR(32) NOT NULL,
#     `registeredAt` DATETIME    NOT NULL,
#     `lastLogin`    DATETIME    NULL DEFAULT NULL,
#     `intro`        TINYTEXT    NULL DEFAULT NULL,
#     `profile`      TEXT        NULL DEFAULT NULL,
#     PRIMARY KEY (`id`),
#     UNIQUE INDEX `uq_mobile` (`mobile` ASC),
#     UNIQUE INDEX `uq_email` (`email` ASC)
# );
#
# CREATE TABLE `user`
# (
#     `id`           bigint(20)  NOT NULL AUTO_INCREMENT,
#     `roleId`       bigint(20)  NOT NULL,
#     `firstName`    varchar(50) DEFAULT NULL,
#     `middleName`   varchar(50) DEFAULT NULL,
#     `lastName`     varchar(50) DEFAULT NULL,
#     `mobile`       varchar(15) DEFAULT NULL,
#     `email`        varchar(50) DEFAULT NULL,
#     `passwordHash` varchar(32) NOT NULL,
#     `registeredAt` datetime    NOT NULL,
#     `lastLogin`    datetime    DEFAULT NULL,
#     `intro`        tinytext    DEFAULT NULL,
#     `profile`      text        DEFAULT NULL,
#     PRIMARY KEY (`id`),
#     UNIQUE KEY `uq_mobile` (`mobile`),
#     UNIQUE KEY `uq_email` (`email`),
#     KEY `idx_user_role` (`roleId`),
#     CONSTRAINT `fk_user_role` FOREIGN KEY (`roleId`) REFERENCES `role` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
# ) ENGINE = InnoDB
#   DEFAULT CHARSET = utf8mb4


#
# use rabbit_gather;
#
# CREATE TABLE `user`
# (
#     `id`                     int unsigned primary key auto_increment,
#     `name`                   varchar(24) not null unique,
#     `password`               char(60)    not null,
#     `randomSalt`             char(24)    not null,
#     `create_time`            timestamp NULL DEFAULT CURRENT_TIMESTAMP,
#     `update_time`            timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
#     `api_permission_bitmask` int unsigned not null
# );
#
# CREATE TABLE `user_information`
# (
#     `user`  int unsigned primary key,
#     `email` varchar(254) unique,
#     `phone` varchar(30) unique,
#     foreign key (`user`) references user (`id`)
# );
#
#
#
CREATE TABLE `article`
(
    `id`          int unsigned primary key auto_increment,
    `title`       varchar(48) not null COMMENT'標題',
    `content`     mediumtext  not null COMMENT'內容',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) COMMENT '文章本體'
;

CREATE TABLE `article_details`
(
    `article` int unsigned primary key,
    `coords`  point not null COMMENT'當前文章釘駐座標',
    foreign key (`article`) references article (`id`)
) COMMENT '文章附屬資訊';

create table `tag_type`
(
    `id`   int unsigned primary key auto_increment,
    `name` char(24) not null unique
) COMMENT '文章標籤的種類';

insert into `tag_type` (name)
    value ('SYSTEM')
;

create table `tags`
(
    `id`   int unsigned primary key auto_increment,
    `name` char(24) not null unique,
    `type` int unsigned  null unique,
        foreign key (`type`) references `tag_type` (`id`)

) COMMENT '文章標籤';


# 代表刪除
insert into `tags` (name,type)
values ('DELETE',1)
;

create table `article_tag`
(
    `tag_id`     int unsigned not null ,
    `article_id` int unsigned not null,
    primary key (`article_id`,`tag_id`),
    foreign key (`article_id`) references `article` (`id`),
    foreign key (`tag_id`) references `tags` (`id`)
) COMMENT '標籤與文章對應';

#
#
#
#

