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

