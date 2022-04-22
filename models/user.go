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
		Username     string              `json:"username" db:"username"`    // Username of the user
		Email        string              `json:"email" db:"email"`          // Email of the user
		DiscordID    nulltype.NullString `json:"discordID" db:"discord_id"` // Discord ID of the user
		PasswordHash string              `json:"-" db:"pwd_hash"`           // Password hash of the user
		Salt         string              `json:"-" db:"salt"`               // Salt used on the password
		AccessLvl    int                 `json:"accessLvl" db:"access_lvl"` // Access level to the api of the user
		CreatedAt    time.Time           `json:"createdAt" db:"created_at"` // Date the user was created
		Banned       bool                `json:"banned" db:"banned"`        // Whether the user is banned or not
	}

	UserCreation struct {
		Username  string              `json:"username" db:"username" form:"username"`     // Username of the user
		Email     string              `json:"email" db:"email" form:"email"`              // Email of the user
		DiscordID nulltype.NullString `json:"discordID" db:"discord_id" form:"discordID"` // Discord ID of the user
		Password  string              `json:"password" db:"-" form:"password"`            // Password of the user
	}

	UserModification struct {
		Email       string              `json:"email" db:"email"`          // Email of the user
		DiscordID   nulltype.NullString `json:"discordID" db:"discord_id"` // Discord ID of the user
		OldPassword string              `json:"oldPassword" db:"-"`        // Old password of the user
		Password    string              `json:"password" db:"-"`           // New password of the user
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
