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
)