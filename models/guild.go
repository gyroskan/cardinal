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
			guild_name=:guild_name, prefix=:prefix, report_channel=:report_channel,
			welcome_channel=:welcome_channel, welcome_message=:welcome_message, 
			private_welcome_msg=:private_welcome_msg, level_channel=:level_channel, level_replace=:level_replace,
			level_response=:level_response,disabled_commands=:disabled_commands,
			allow_moderation=:allow_moderation, max_warns=:max_warns, ban_time=:ban_time
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
		GuildID           string              `json:"guildID" db:"guild_id"`                      // Guild ID
		GuildName         string              `json:"guildName" db:"guild_name"`                  // Name of the guild
		Prefix            string              `json:"prefix" db:"prefix"`                         // Prefix used for calling the bot
		ReportChannel     nulltype.NullString `json:"reportChannel" db:"report_channel"`          // Channel ID for reporting
		WelcomeChannel    nulltype.NullString `json:"welcomeChannel" db:"welcome_channel"`        // Channel ID to send welcome messages
		WelcomeMsg        nulltype.NullString `json:"welcomeMsg" db:"welcome_message"`            // Message to send when a user joins
		PrivateWelcomeMsg nulltype.NullString `json:"privateWelcomeMsg" db:"private_welcome_msg"` // Message to send when a user joins in DM
		LvlChannel        nulltype.NullString `json:"lvlChannel" db:"level_channel"`              // Channel ID to send level up messages
		LvlReplace        bool                `json:"lvlReplace" db:"level_replace"`              // Weather or not to replace previous rewards
		LvlResponse       int                 `json:"lvlResponse" db:"level_response"`            // If the level is a multiple of this number, send a level up message
		DisabledCommands  nulltype.NullString `json:"disabledCommands" db:"disabled_commands"`    // List of disabled commands separated by slashes
		AllowModeration   bool                `json:"allowModeration" db:"allow_moderation"`      // Whether or not to allow moderation commands
		MaxWarns          int                 `json:"maxWarns" db:"max_warns"`                    // Max number of warnings before a user is banned
		BanTime           int                 `json:"banTime" db:"ban_time"`                      // Time in days to ban a user for
		// TODO is Members field needed?
	}

	GuildPres struct {
		GuildID   string `json:"guildID" db:"guildID"`     // Guild ID
		GuildName string `json:"guildName" db:"guildname"` // Name of the guild
		Prefix    string `json:"prefix" db:"prefix"`       // Prefix used for calling the bot
	}
)
