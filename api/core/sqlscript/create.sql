use core;
drop table if exists wallet_log;
drop table if exists wallet_credit;
drop table if exists transfer_log;
drop table if exists transfer_order;
drop table if exists transfer_target;
drop table if exists wallet;
drop table if exists user_role;
drop table if exists user_info;
drop table if exists user;
drop table if exists restful_rbac_model;
drop table if exists role;
drop table if exists currency;

drop table if exists joy_games_baccarat_order_transfer_order_map;
drop table if exists joy_games_baccarat_round_order_map;
drop table if exists project_card_list;



CREATE TABLE `restful_rbac_model`
(
    `p_type` varchar(32)  NOT NULL DEFAULT '',
    `v0`     varchar(255) NOT NULL DEFAULT '',
    `v1`     varchar(255) NOT NULL DEFAULT '',
    `v2`     varchar(255) NOT NULL DEFAULT '',
    `v3`     varchar(255) NOT NULL DEFAULT '',
    `v4`     varchar(255) NOT NULL DEFAULT '',
    `v5`     varchar(255) NOT NULL DEFAULT '',
    KEY `idx_restful_rbac_model` (`p_type`, `v0`, `v1`),
    unique (`p_type`, `v0`, `v1`, `v2`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


CREATE TABLE `role`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at`  TIMESTAMP                    DEFAULT CURRENT_TIMESTAMP,
    `name`        varchar(255)        NOT NULL,
    `description` tinytext                     DEFAULT NULL,
    `active`      tinyint(1)          NOT NULL default 1,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`)

)
    ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4;



CREATE TABLE `user`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `type`          char(10)            not null,
    `uuid`          char(20)            NOT NULL,
    `created_at`    TIMESTAMP                    DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    TIMESTAMP                    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`    TIMESTAMP                    DEFAULT 0,
    `account`       varchar(150)        NOT NULL,
    `password_hash` char(60)            NOT NULL,
    `password_salt` char(24)            NOT NULL,
    `last_login`    datetime                     DEFAULT NULL,
    `frozen`        tinyint(1)          NOT NULL DEFAULT 0,
    UNIQUE KEY (`id`),
    UNIQUE KEY (`uuid`),
    PRIMARY KEY (`id`, `uuid`),
    UNIQUE KEY (`account`),
    KEY (`deleted_at`)
) auto_increment = 10000
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


CREATE TABLE `user_info`
(
    `user_id`   bigint(20) unsigned NOT NULL,
    `nick_name` varchar(128)        not null comment '暱稱(一開始與帳號一樣)',
    `gender`    tinyint(1)          not null comment '性別(女0男1不透露3)',
    `email`     varchar(255)        not null,
    `phone`     varchar(255)  default '',
    `photo_url` VARCHAR(2083) default '',
    foreign key (`user_id`) references user (`id`),
    unique (`email`)
) auto_increment = 10000
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE `user_role`
(
    `user_id`    bigint(20) unsigned NOT NULL,
    `role_id`    bigint(20) unsigned NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`, `role_id`),
    KEY (`role_id`),
    FOREIGN KEY (`role_id`) REFERENCES `role` (`id`) ON UPDATE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON UPDATE CASCADE
)
    ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4;



# 需做存讀取暫存
create table `project_card_list`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `title`      varchar(1024)       not null comment '標題',
    `summary`    varchar(1024)       not null comment '簡介',
    `content_url` VARCHAR(2083) not null  comment '內容網址',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP comment '創建時間',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新時間',
    `deleted_at` TIMESTAMP DEFAULT 0 comment '刪除時間',
    PRIMARY KEY (`id`),
    KEY (`created_at`),
    KEY (`updated_at`),
    KEY (`deleted_at`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4;



insert into role (name) value ('login');

insert into project_card_list (title, summary, content_url) value ( concat('title' , LAST_INSERT_ID()+1) , concat('summary' ,LAST_INSERT_ID()+1) ,concat('content_url' , LAST_INSERT_ID()+1) );
insert into project_card_list (title, summary, content_url) value ( concat('title' , LAST_INSERT_ID()+1) , concat('summary' ,LAST_INSERT_ID()+1) ,concat('content_url' , LAST_INSERT_ID()+1) );
insert into project_card_list (title, summary, content_url) value ( concat('title' , LAST_INSERT_ID()+1) , concat('summary' ,LAST_INSERT_ID()+1) ,concat('content_url' , LAST_INSERT_ID()+1) );
insert into project_card_list (title, summary, content_url) value ( concat('title' , LAST_INSERT_ID()+1) , concat('summary' ,LAST_INSERT_ID()+1) ,concat('content_url' , LAST_INSERT_ID()+1) );
insert into project_card_list (title, summary, content_url) value ( concat('title' , LAST_INSERT_ID()+1) , concat('summary' ,LAST_INSERT_ID()+1) ,concat('content_url' , LAST_INSERT_ID()+1) );
insert into project_card_list (title, summary, content_url) value ( concat('title' , LAST_INSERT_ID()+1) , concat('summary' ,LAST_INSERT_ID()+1) ,concat('content_url' , LAST_INSERT_ID()+1) );
insert into project_card_list (title, summary, content_url) value ( concat('title' , LAST_INSERT_ID()+1) , concat('summary' ,LAST_INSERT_ID()+1) ,concat('content_url' , LAST_INSERT_ID()+1) );
insert into project_card_list (title, summary, content_url) value ( concat('title' , LAST_INSERT_ID()+1) , concat('summary' ,LAST_INSERT_ID()+1) ,concat('content_url' , LAST_INSERT_ID()+1) );

select * from project_card_list;


