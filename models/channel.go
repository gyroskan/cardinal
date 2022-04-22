package models

const (
	CreateChannelQuery = `
		INSERT INTO channel
			(channel_id, guild_id, ignored, xp_blacklisted)
		VALUES
			(:channel_id, :guild_id, :ignored, :xp_blacklisted)
	`
	UpdateChannelQuery = `
		UPDATE channel SET
			ignored=:ignored, xp_blacklisted=:xp_blacklisted
		WHERE
			guild_id=:guild_id AND channel_id=:channel_id
	`
)

type (
	Channel struct {
		ChannelID     string `json:"channelID" db:"channel_id"`         // ID of the channel
		GuildID       string `json:"guildID" db:"guild_id"`             // ID of the guild
		Ignored       bool   `json:"ignored" db:"ignored"`              // Wether the channel is ignored by the bot or not
		XpBlacklisted bool   `json:"wpBlacklisted" db:"xp_blacklisted"` // Wether the channel is blacklisted from xp or not
	}
)
