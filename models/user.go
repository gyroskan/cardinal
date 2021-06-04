package models

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/mattn/go-nulltype"
)

/* const */
var discRegex = regexp.MustCompile(`^[0-9]{17,21}`)

const (
	InsertUserQuery = `
		INSERT INTO user
			(username, email, discord_id, pwd_hash, salt, access_lvl, created_at, banned)
		VALUES
			(:username,:email,:discord_id,:pwd_hash,:salt,:access_lvl,:created_at,:banned)
			`
	UpdateUserQuery = `
		UPDATE user SET
			email=:email, discord_id=:discord_id, pwd_hash=:pwd_hash, salt=:salt
		WHERE 
			username=:username
			`
)

type (
	User struct {
		Username     string              `json:"username" db:"username"`
		Email        string              `json:"email" db:"email"`
		DiscordID    nulltype.NullString `json:"discordID" db:"discord_id"`
		PasswordHash string              `json:"-" db:"pwd_hash"`
		Salt         string              `json:"-" db:"salt"`
		AccessLvl    int                 `json:"accessLvl" db:"access_lvl"`
		CreatedAt    time.Time           `json:"createdAt" db:"created_at"`
		Banned       bool                `json:"banned" db:"banned"`
	}

	UserCreation struct {
		Username  string              `json:"username" db:"username" form:"username"`
		Email     string              `json:"email" db:"email" form:"email"`
		DiscordID nulltype.NullString `json:"discordID" db:"discord_id" form:"discordID"`
		Password  string              `json:"password" db:"-" form:"password"`
	}

	UserModification struct {
		Email       string              `json:"email" db:"email"`
		DiscordID   nulltype.NullString `json:"discordID" db:"discord_id"`
		OldPassword string              `json:"oldPassword" db:"-"`
		Password    string              `json:"password" db:"-"`
	}
)

// Validate UserCreation fields
func (u *UserCreation) Validate() error {
	u.Username = strings.TrimSpace(u.Username)
	if u.Username == "" {
		return errors.New("username empty")
	}
	if u.Username != strings.ToLower(u.Username) {
		return errors.New("username must be lowercase")
	}
	if l := len(u.Username); l < 3 || l > 20 {
		return errors.New("username too long")
	}

	if u.Email == "" {
		return errors.New("email empty")
	}
	if l := len(u.Email); l > 80 {
		return errors.New("email too long")
	}
	m, err := mail.ParseAddress(u.Email)
	if err != nil {
		return errors.New("email invalid")
	}
	u.Email = m.Address

	if u.DiscordID.Valid() && !discRegex.MatchString(u.DiscordID.StringValue()) {
		return errors.New("invalid discordID")
	}

	if u.Password == "" {
		return errors.New("password empty")
	}

	return nil
}
