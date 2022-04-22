package models

import (
	"time"

	"github.com/mattn/go-nulltype"
)

const (
	CreateBanQuery = `
		INSERT INTO ban 
			(member_id, guild_id, banner_id, banned_at, ban_reason, auto_ban)
		VALUES
			(:member_id, :guild_id, :banner_id, :banned_at, :ban_reason, :auto_ban)
	`
)

type (
	Ban struct {
		BanID     int                 `json:"banID" db:"ban_id"`         // ID of the ban
		MemberID  string              `json:"memberID" db:"member_id"`   // ID of the member
		GuildID   string              `json:"guildID" db:"guild_id"`     // ID of the guild
		BannerID  nulltype.NullString `json:"bannerID" db:"banner_id"`   // ID of the user who banned the member
		BannedAt  time.Time           `json:"bannedAt" db:"banned_at"`   // Date the member was banned
		BanReason nulltype.NullString `json:"banReason" db:"ban_reason"` // Reason for the ban
		AutoBan   bool                `json:"autoBan" db:"auto_ban"`     // Whether the ban was automatic or not
	}
)
