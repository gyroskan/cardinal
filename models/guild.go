package models

import "github.com/mattn/go-nulltype"

const (
	CreateGuildQuery = `
		INSERT INTO guild
			(guild_id, guild_name, prefix, report_channel, welcome_channel, welcome_message,
			private_welcome_msg, level_channel, lvl_replace, level_response, disabled_commands,
			allow_moderation, max_warns, ban_time)
		VALUES
			(:guild_id, :guild_name, :prefix, :report_channel, :welcome_channel, :welcome_message,
			:private_welcome_msg, :level_channel, :lvl_replace, :level_response, :disabled_commands,
			:allow_moderation, :max_warns, :ban_time)
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
		LvlReplace        bool                `json:"lvlReplace" db:"lvl_replace"`
		LvlResponse       int                 `json:"lvlResponse" db:"level_response"`
		DisabledCommands  nulltype.NullString `json:"disabledCommands" db:"disabled_commands"`
		AllowModeration   bool                `json:"allowModeration" db:"allow_moderation"`
		MaxWarns          int                 `json:"maxWarns" db:"max_warns"`
		BanTime           int                 `json:"banTime" db:"ban_time"`
		Members           []Member            `json:"members" db:"-"`
	}

	GuildPres struct {
		GuildID   string `json:"guildID" db:"guildID"`
		GuildName string `json:"guildName" db:"guildname"`
		Prefix    string `json:"prefix" db:"prefix"`
	}
)
