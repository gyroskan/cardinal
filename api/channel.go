package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gyroskan/cardinal/db"
	"github.com/gyroskan/cardinal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func initChannels() {
	chans := apiGroupe.Group("/guilds/:guildID/chans")
	chans.GET("/", getChannels).Name = "Fetch channels of a guild."
	chans.GET("/:id", getChannel).Name = "Fetch channel of a guild."
	chans.POST("/", createChannel).Name = "Create channel."
	chans.PATCH("/:id", updateChannel).Name = "Update channel values."
	chans.DELETE("/:id", deleteChannel).Name = "Delete channel."
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

func getChannel(c echo.Context) error {
	guildID := c.Param("guildID")
	chanID := c.Param("id")

	var chann models.Channel

	err := db.DB.Get(&chann, "SELECT * FROM channel WHERE guild_id=?,channel_id=?", guildID, chanID)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, nil)
		}
		log.Error("getChannel/ Error retrieving channel: ", err)
		c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, chann)
}

func createChannel(c echo.Context) error {
	var channel models.Channel
	guildID := c.Param("guildID")

	if err := c.Bind(&channel); err != nil || channel.GuildID != guildID {
		return echo.NewHTTPError(http.StatusBadRequest, channel)
	}

	_, err := db.DB.NamedExec(models.CreateChannelQuery, channel)

	if err != nil {
		// TODO switch case of sql errors.
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, channel)
}

func updateChannel(c echo.Context) error {
	guildID := c.Param("guildID")
	chanID := c.Param("id")

	var channel models.Channel

	err := db.DB.Get(&channel, "SELECT * FROM channel WHERE guild_id=?,channel_id=?", guildID, chanID)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, nil)
		}
		log.Error("getChannel/ Error retrieving channel: ", err)
		c.JSON(http.StatusInternalServerError, nil)
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&channel); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	channel.ChannelID = chanID
	channel.GuildID = guildID

	_, err = db.DB.NamedExec(models.UpdateChannelQuery, channel)
	if err != nil {
		log.Warn("UpdateMember/ Error updating member: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, channel)
}

func deleteChannel(c echo.Context) error {
	guildID := c.Param("guildID")
	chanID := c.Param("id")

	res, err := db.DB.Exec("DELETE FROM channel WHERE guild_id = ?,channel_id = ?", guildID, chanID)

	if err != nil {
		log.Warn("HardDeleteMember/ Error while deleting member from db: ", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not delete the channel."})
	}

	if r, _ := res.RowsAffected(); r == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}
