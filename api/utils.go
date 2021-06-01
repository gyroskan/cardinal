package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func generateSalt() ([]byte, error) {
	var salt = make([]byte, 16)

	_, err := rand.Read(salt[:])

	if err != nil {
		log.Warn("generateSalt/:", err)
		return nil, err
	}

	return salt, nil
}

func hashPassword(password string, salt []byte) string {
	var passwordBytes = []byte(password)
	passwordBytes = append(passwordBytes, salt...)

	var hasher = sha256.New()
	hasher.Write(passwordBytes)

	return hex.EncodeToString(hasher.Sum(nil))
}

func isAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*JwtCustomClaims)
		accessLevel := claims.Access_level
		if accessLevel != 0 {
			return echo.ErrForbidden
		}
		return next(c)
	}
}

func isAdminOrLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*JwtCustomClaims)
		accessLevel := claims.Access_level
		if accessLevel != 0 && c.Param("username") != claims.Username {
			return echo.ErrForbidden
		}
		return next(c)
	}
}
