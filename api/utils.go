package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

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
