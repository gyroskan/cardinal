package models

const (
	CreateRoleQuery = `
		INSERT INTO role
			(role_id, guild_id, is_default, ignored, reward, xp_blacklisted)
		VALUES
			(:role_id, :guild_id, :is_default, :ignored, :reward, :xp_blacklisted)
	`
)

type (
	Role struct {
		RoleID        string `json:"roleID" db:"role_id"`
		GuildID       string `json:"guildID" db:"guild_id"`
		IsDefault     bool   `json:"isDefault" db:"is_default"`
		Reward        int    `json:"reward" db:"reward"`
		Ignored       bool   `json:"ignored" db:"ignored"`
		XpBlacklisted bool   `json:"wpBlacklisted" db:"xp_blacklisted"`
	}
)
