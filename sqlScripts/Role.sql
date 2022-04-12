CREATE TABLE IF NOT EXISTS `role`
(
    `role_id`        varchar(20) not null,               -- PK ID
    `guild_id`       varchar(20) not null,               -- FK guild_PK
    `is_default`     tinyint(1)  not null default false, -- Whether to give the role to new member or not
    `ignored`        tinyint(1)  not null default false, -- If true, bot does not respond to user with this role
    `reward`         int         not null default 0,     -- The level corresponding to the reward
    `xp_blacklisted` tinyint(1)  not null default false, -- No xp leveling for this role
    constraint `fk_role_guild_id`
        foreign key (guild_id)
            references guild (guild_id)
            on delete cascade,
    constraint `pk_channel`
        primary key (role_id, guild_id)
)