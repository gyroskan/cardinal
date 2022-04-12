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