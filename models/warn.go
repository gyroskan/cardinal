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
		WarnID     int                 `json:"warnID" db:"warn_id"`         // ID of the warn
		MemberID   string              `json:"memberID" db:"member_id"`     // ID of the member
		GuildID    string              `json:"guildID" db:"guild_id"`       // ID of the guild
		WarnerID   nulltype.NullString `json:"warnerID" db:"warner_id"`     // ID of the user who warned the member
		WarnedAt   time.Time           `json:"warnedAt" db:"warned_at"`     // Date the member was warned
		WarnReason nulltype.NullString `json:"warnReason" db:"warn_reason"` // Reason for the warn
	}
)
