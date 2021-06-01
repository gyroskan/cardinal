package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	version = "v1"
)

var (
	apiGroupe *echo.Group
)

func InitRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	apiGroupe = e.Group("/api/" + version)

	// swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "swagger/index.html")
	})

	initAuth()

	config := middleware.JWTConfig{
		Claims:     &JwtCustomClaims{},
		SigningKey: []byte(secret),
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
	fmt.Println("Started cardinal API" + version + ", made by gyroskan!")
}
