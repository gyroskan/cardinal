package models

import "github.com/mattn/go-nulltype"

const (
	SelectGuildMembersQuery = `
			SELECT * FROM member
			WHERE guild_id=? AND member_id > ?
			ORDER BY member_id ASC
			LIMIT ?
		`
	CreateMemberQuery = `
		INSERT INTO member 
			(member_id, guild_id, joined_at, ` + "`left`" + `, xp, level)
		VALUES
			(:member_id, :guild_id, :joined_at, :left, :xp, :level)
		`
	ResetMemberQuery = `
		UPDATE member SET
			` + "`left`" + `=DEFAULT, xp=DEFAULT, level=DEFAULT
		WHERE
			guild_id=? AND member_id=?
		`
	ResetGuildMembersQuery = `
		UPDATE member SET
			` + "`left`" + `=DEFAULT, xp=DEFAULT, level=DEFAULT
		WHERE
			guild_id=?
		`
	UpdateMemberQuery = `
		UPDATE member SET 
			` + "`left`" + `=:left, xp=:xp, level=:level
		WHERE 
			guild_id=:guild_id AND member_id=:member_id
		`
)

type Member struct {
	MemberID string            `json:"memberID" db:"member_id"`                    // Member ID
	GuildID  string            `json:"guildID" db:"guild_id"`                      // Guild ID
	JoinedAt nulltype.NullTime `json:"joinedAt" db:"joined_at" format:"date-time"` // Date for when the member joined the guild
	Left     int               `json:"left" db:"left"`                             // Number of times the member left the guild
	Xp       int               `json:"xp" db:"xp"`                                 // Amount of xp the member has
	Level    int               `json:"level" db:"level"`                           // Level of the member
}
