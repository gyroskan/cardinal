package models

import (
	"time"

	"github.com/mattn/go-nulltype"
)

const (
	CreateWarnQuery = `
		INSERT INTO warn 
			(member_id, guild_id, warner_id, warned_at, warn_reason)
		VALUES
			(:member_id, :guild_id, :warner_id, :warned_at, :warn_reason)
	`
)

type (
	Warn struct {
		WarnID     int                 `json:"warnID" db:"warn_id"`
		MemberID   string              `json:"memberID" db:"member_id"`
		GuildID    string              `json:"guildID" db:"guild_id"`
		WarnerID   nulltype.NullString `json:"warnerID" db:"warner_id"`
		WarnedAt   time.Time           `json:"warnedAt" db:"warned_at"`
		WarnReason nulltype.NullString `json:"warnReason" db:"warn_reason"`
	}
)
