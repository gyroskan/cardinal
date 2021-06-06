package models

const ()

type (
	Channel struct {
		ChannelID     string `json:"channelID" db:"channel_id"`
		GuildID       string `json:"guildID" db:"guild_id"`
		Ignored       bool   `json:"ignored" db:"ignored"`
		XpBlacklisted bool   `json:"wpBlacklisted" db:"xp_blacklisted"`
	}
)
