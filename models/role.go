package models

const (
	CreateRoleQuery = `
		INSERT INTO role
			(role_id, guild_id, is_default, ignored, reward, xp_blacklisted)
		VALUES
			(:role_id, :guild_id, :is_default, :ignored, :reward, :xp_blacklisted)
	`
	UpdateRoleQuery = `
		UPDATE role SET
			is_default=:is_default, ignored=:ignored, reward=:reward, xp_blacklisted=:xp_blacklisted
		WHERE
			guild_id=:guild_id AND role_id=:role_id
	`
)

type (
	Role struct {
		RoleID        string `json:"roleID" db:"role_id"`               // ID of the role
		GuildID       string `json:"guildID" db:"guild_id"`             // ID of the guild
		IsDefault     bool   `json:"isDefault" db:"is_default"`         // Wether to give the role to new members
		Reward        int    `json:"reward" db:"reward"`                // The level corresponding to the reward
		Ignored       bool   `json:"ignored" db:"ignored"`              // Wether the role is ignored by the bot or not
		XpBlacklisted bool   `json:"xpBlacklisted" db:"xp_blacklisted"` // Wether the role is blacklisted from xp or not
	}
)
