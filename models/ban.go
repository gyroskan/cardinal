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
		BanID     int                 `json:"banID" db:"ban_id"`
		MemberID  string              `json:"memberID" db:"member_id"`
		GuildID   string              `json:"guildID" db:"guild_id"`
		BannerID  nulltype.NullString `json:"bannerID" db:"banner_id"`
		BannedAt  time.Time           `json:"bannedAt" db:"banned_at"`
		BanReason nulltype.NullString `json:"banReason" db:"ban_reason"`
		AutoBan   bool                `json:"autoBan" db:"auto_ban"`
	}
)
