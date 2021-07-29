package models

import "github.com/mattn/go-nulltype"

const (
	CreateGuildQuery = `
		INSERT INTO guild
			(guild_id, guild_name, prefix, report_channel, welcome_channel, welcome_message,
			private_welcome_msg, level_channel, level_replace, level_response, disabled_commands,
			allow_moderation, max_warns, ban_time)
		VALUES
			(:guild_id, :guild_name, :prefix, :report_channel, :welcome_channel, :welcome_message,
			:private_welcome_msg, :level_channel, :level_replace, :level_response, :disabled_commands,
			:allow_moderation, :max_warns, :ban_time)
		`
	UpdateGuildQuery = `
		UPDATE guild SET
			(guild_name=:guild_name, prefix=:prefix, report_channel=:report_channel,
			welcome_channel=:welcome_channel, welcome_message=:welcome_message, 
			private_welcome_msg=:private_welcome_msg, level_channel=:level_channel, level_replace=:level_replace,
			level_response=:level_response,disabled_commands=:disabled_commands,
			allow_moderation=:allow_moderation, max_warns=:max_warns, ban_time=:ban_time)
		WHERE
			guild_id=:guild_id
		`
	ResetGuildQuery = `
		UPDATE guild SET
			prefix=DEFAULT,report_channel=DEFAUT,welcome_channel=DEFAUT, welcome_message=DEFAULT,
			private_welcome_msg=DEFAULT,level_channel=DEFAUT,level_response=DEFAULT,level_replace=DEFAULT,
			allow_moderation=DEFAULT, max_warns=DEFAULT, ban_time=DEFAULT
		WHERE
			guild_id=?
	`
)

type (
	Guild struct {
		GuildID           string              `json:"guildID" db:"guild_id"`
		GuildName         string              `json:"guildName" db:"guild_name"`
		Prefix            string              `json:"prefix" db:"prefix"`
		ReportChannel     nulltype.NullString `json:"reportChannel" db:"report_channel"`
		WelcomeChannel    nulltype.NullString `json:"welcomeChannel" db:"welcome_channel"`
		WelcomeMsg        nulltype.NullString `json:"welcomeMsg" db:"welcome_message"`
		PrivateWelcomeMsg nulltype.NullString `json:"privateWelcomeMsg" db:"private_welcome_msg"`
		LvlChannel        nulltype.NullString `json:"lvlChannel" db:"level_channel"`
		LvlReplace        bool                `json:"lvlReplace" db:"level_replace"`
		LvlResponse       int                 `json:"lvlResponse" db:"level_response"`
		DisabledCommands  nulltype.NullString `json:"disabledCommands" db:"disabled_commands"`
		AllowModeration   bool                `json:"allowModeration" db:"allow_moderation"`
		MaxWarns          int                 `json:"maxWarns" db:"max_warns"`
		BanTime           int                 `json:"banTime" db:"ban_time"`
		// TODO is Members field needed?
	}

	GuildPres struct {
		GuildID   string `json:"guildID" db:"guildID"`
		GuildName string `json:"guildName" db:"guildname"`
		Prefix    string `json:"prefix" db:"prefix"`
	}
)
