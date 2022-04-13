DROP DATABASE IF EXISTS `cardinal`;
CREATE DATABASE `cardinal`;
USE `cardinal`;

create table if not exists `guild`
(
    `guild_id`            varchar(20)                       not null, -- PK ID
    `guild_name`          varchar(102)                      not null, -- Name of the guild
    `prefix`              varchar(8)          default '%'   not null, -- Prefix used to call the bot
    `report_channel`      varchar(20)         default null  null,     -- Channel id for report
    `welcome_channel`     varchar(20)         default null  null,     -- Channel to send the welcome_message
    `welcome_message`     text                default null  null,     -- Message sent to new members
    `private_welcome_msg` text                default null  null,     -- Message sent in DM to new member
    `level_channel`       varchar(20)         default null  null,     -- If specified, channel for level up messaged
    `level_response`      tinyint(1) unsigned default 1     not null, -- % Level when the bot respond for level up
    `level_replace`       tinyint(1)          default false not null, -- Whether to replace role rewards or not
    `disabled_commands`   text                              null,     -- List of disabled commands separated by '/'
    `allow_moderation`    tinyint(1)          default true  not null, -- Whether moderation's command are allowed or not
    `max_warns`           tinyint(1) unsigned default 3     not null, -- Max warnings before ban of the member
    `ban_time`            tinyint(1) unsigned default 7     not null, -- Time of an automatic ban
    constraint `pk_guild`
        primary key (guild_id)
);

CREATE TABLE IF NOT EXISTS `member`
(
    `member_id` varchar(20)          NOT NULL,                             -- Discord ID of the user.
    `guild_id`  varchar(20)          NOT NULL,                             -- Discord ID of the guild.
    `joined_at` datetime                      DEFAULT current_timestamp(), -- Date when user joined the guild.
    `left`      tinyint(3) unsigned  NOT NULL DEFAULT 0,                   -- Number of times that user left the guild.
    `xp`        int(10) unsigned     NOT NULL DEFAULT 0,                   -- Amount of xp of the user.
    `level`     smallint(5) unsigned NOT NULL DEFAULT 0,                   -- Level of the user
    constraint `fk_member_guild_id`
        foreign key (guild_id)
            references guild (guild_id)
            on delete cascade,
    constraint `pk_member`
        PRIMARY KEY (member_id, guild_id)                                  -- User can be in multiple guild.
);

CREATE TABLE IF NOT EXISTS `user`
(
    `username`   varchar(20) NOT NULL,                             -- Username of the user
    `email`      varchar(80) NOT NULL,                             -- Email of the user.
    `discord_id` varchar(20),                                      -- Discord ID of the user.
    `pwd_hash`   VARCHAR(64) NOT NULL,                             -- Password hash.
    `salt`       VARCHAR(32) NOT NULL,                             -- Salt used to hash password.
    `access_lvl` int         NOT NULL,                             -- Access level of the user. (0 admin, 1 all except api users, 2 get only)
    `created_at` DATETIME    NOT NULL default CURRENT_TIMESTAMP(), -- Date of creation.
    `banned`     BOOLEAN     NOT NULL,                             -- Whether the user is banned.
    constraint `pk_user`
        PRIMARY KEY (username)
);

create table if not exists `channel`
(
    `channel_id`        varchar(20) not null,               -- PK ID
    `guild_id`          varchar(20) not null,               -- FK guild_PK
    `ignored`           tinyint(1)  not null default false, -- If true, the bot does not respond in this channel
    `xp_blacklisted`    tinyint(1)  not null default false, -- If true, no xp gained in this channels
    constraint `fk_channel_guild_id`
        foreign key (guild_id)
            references guild (guild_id)
            on delete cascade,
    constraint `pk_channel`
        primary key (channel_id, guild_id)
);

CREATE TABLE IF NOT EXISTS `ban`
(
    `ban_id`     int         not null auto_increment,              -- PK ID
    `member_id`  varchar(20) not null,                             -- Member id of the banned user.
    `guild_id`   varchar(20) not null,                             -- Guild id from which the user was banned.
    `banner_id`  varchar(20),                                      -- Id of the user who banned.
    `banned_at`  datetime    not null default current_timestamp(), -- Date of the ban.
    `ban_reason` text,                                             -- Reason of the ban
    `auto_ban`   tinyint(1)  not null default false,               -- Whether it was a banned after multiple warning or not.
    constraint `fk_ban_member_pk`
        foreign key (member_id, guild_id)
            references member (member_id, guild_id)
            on delete cascade,
    constraint `pk_ban`
        primary key (ban_id)
);

CREATE TABLE IF NOT EXISTS `warn`
(
    `warn_id`     int         NOT NULL AUTO_INCREMENT,              -- PK ID
    `member_id`   varchar(20) NOT NULL,                             -- Discord ID of the user.
    `guild_id`    varchar(20) NOT NULL,                             -- Guild ID where the warn happened.
    `warner_id`   varchar(20),                                      -- Discord ID of the member who created the warn.
    `warned_at`   datetime    NOT NULL default current_timestamp(), -- Date of the creation of the warn.
    `warn_reason` text,                                             -- Reason of the warn.
    constraint `fk_warn_member_pk`
        foreign key (member_id, guild_id)
            references member (member_id, guild_id)
            on delete cascade,
    constraint `pk_warn`
        primary key (warn_id)
);
