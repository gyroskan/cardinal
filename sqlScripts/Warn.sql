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