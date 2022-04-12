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
)