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
			guild_id=:guild_id,channel_id=:channel_id
	`
)

type (
	Channel struct {
		ChannelID     string `json:"channelID" db:"channel_id"`
		GuildID       string `json:"guildID" db:"guild_id"`
		Ignored       bool   `json:"ignored" db:"ignored"`
		XpBlacklisted bool   `json:"wpBlacklisted" db:"xp_blacklisted"`
	}
)
