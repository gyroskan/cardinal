package api

import (
	"net/http"

	"github.com/gyroskan/cardinal/db"
	"github.com/gyroskan/cardinal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func initChannels() {
	chans := apiGroupe.Group("/guilds/:guildID/chans")
	chans.GET("/", getChannels).Name = "Fetch channels of a guild."
}

func getChannels(c echo.Context) error {
	guildID := c.Param("guildID")
	var channels models.Channel

	err := db.DB.Select(&channels, "SELECT * FROM channel WHERE guild_id=?", guildID)

	if err != nil {
		log.Warn("GetChannels/ Error retrieving channels from guildID: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, channels)
}
