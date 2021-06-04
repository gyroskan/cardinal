package api

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/gyroskan/cardinal/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	version   = "v1"
	base_path = "/api/v1"
)

var (
	apiGroupe *echo.Group
)

// @title Cardinal API
// @version 0.1
// @description The API to interact with cardinal discord bot database.

// @contact.name API Support
// @contact.email gyroskan@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host 	cardinal.gyroskan.com:5005
// @BasePath /api/v1
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
			return c.Request().URL.Path == base_path+"/users/register" ||
				c.Request().URL.Path == base_path+"/users/login"
		},
	}
	apiGroupe.Use(middleware.JWTWithConfig(config))

	initUsers()

	return e
}

func Run() {
	e := InitRouter()

	if err := e.Start(":5005"); err != nil {
		log.Fatal("unable to start api. ", err)
	}
	fmt.Println("Started cardinal API " + version + ", made by gyroskan!")
}
