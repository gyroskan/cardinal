package api

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	_ "github.com/gyroskan/cardinal/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	version   = "v1.0.3"
	base_path = "/api/v1"
)

var (
	apiGroupe *echo.Group
)

// @title Cardinal API
// @version 1.0.3
// @description The API to interact with cardinal discord bot database.

// @contact.name API Support
// @contact.email gyroskan@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @Schemes                     https http
// @host                        cardinal.gyroskan.com
// @BasePath                    /api/v1
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
func InitRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	apiGroupe = e.Group(base_path)

	// swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "swagger/index.html")
	})

	initAuth()

	config := middleware.JWTConfig{
		Claims:     &JwtCustomClaims{},
		SigningKey: []byte(secret),
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == base_path+"/users/register" ||
				c.Request().URL.Path == base_path+"/users/login" {
				return true
			}
			if c.Get("user") == nil {
				return false
			}
			user, success := c.Get("user").(*jwt.Token)
			if !success {
				return false
			}
			claims, success := user.Claims.(*JwtCustomClaims)
			if !success {
				return false
			}
			accessLevel := claims.Access_level

			if c.Request().Method == "GET" {
				return accessLevel <= 2
			}

			return strings.HasPrefix(c.Request().URL.Path, "/users") || accessLevel <= 1
		},
	}
	apiGroupe.Use(middleware.JWTWithConfig(config))

	initUsers()
	initGuilds()
	initMembers()
	initChannels()
	initRoles()
	initBans()
	initWarn()

	return e
}

func Run() {
	e := InitRouter()

	log.Info("Started cardinal API " + version + ", made by gyroskan!")
	if err := e.Start(":5005"); err != nil {
		log.Fatal("unable to start api. ", err)
	}
}
