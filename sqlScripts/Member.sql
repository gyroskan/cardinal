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