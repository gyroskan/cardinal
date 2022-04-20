package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gyroskan/cardinal/db"
	"github.com/gyroskan/cardinal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func initChannels() {
	chans := apiGroupe.Group("/guilds/:guildID/channels")
	chans.GET("/", getChannels).Name = "Fetch channels of a guild."
	chans.GET("/:id", getChannel).Name = "Fetch channel of a guild."
	chans.POST("/", createChannel).Name = "Create channel."
	chans.PATCH("/:id", updateChannel).Name = "Update channel values."
	chans.DELETE("/:id", deleteChannel).Name = "Delete channel."
}

// @Summary      Get Guild channels
// @Tags         Channels
// @Description  Fetch all channels of the guild.
// @Param        guildID      path   string   true   "guild id"
// @Param        ignored      query  bool     false  "ignored channels only"      default(false)
// @Param        xpBlacklist  query  bool     false  "xpBlacklist channels only"  default(false)
// @Success      200          "OK"   {array}  models.Channel
// @Failure      403          "Forbidden"
// @Failure      500          "Server error"
// @Router       /guilds/{guildID}/channels [GET]
func getChannels(c echo.Context) error {
	guildID := c.Param("guildID")
	ignored := false
	xpBlacklisted := false
	if c.QueryParam("ignored") != "" {
		ignored, _ = strconv.ParseBool(c.QueryParam("ignored"))
	}
	if c.QueryParam("xpBlacklist") != "" {
		xpBlacklisted, _ = strconv.ParseBool(c.QueryParam("xpBlacklist"))
	}
	var channels []models.Channel

	query := `SELECT * FROM channel WHERE guild_id=?`
	if ignored {
		query += " AND ignored=true"
	}
	if xpBlacklisted {
		query += " AND xp_blacklisted=true"
	}

	err := db.DB.Select(&channels, query, guildID)

	if err != nil {
		log.Warn("GetChannels/ Error retrieving channels from guildID: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, channels)
}

// @Summary      Get one Guild Channel
// @Tags         Channels
// @Description  Fetch the channel of the guild.
// @Param        guildID    path  string    true  "guild id"
// @Param        channelID  path  string    true  "channel id"
// @Success      200        "OK"  {object}  models.Channel
// @Failure      403        "Forbidden"
// @Failure      404        "Not Found"
// @Failure      500        "Server error"
// @Router       /guilds/{guildID}/channels/{channelID} [GET]
func getChannel(c echo.Context) error {
	guildID := c.Param("guildID")
	chanID := c.Param("id")

	var chann models.Channel

	err := db.DB.Get(&chann, "SELECT * FROM channel WHERE guild_id=? AND channel_id=?", guildID, chanID)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, nil)
		}
		log.Error("getChannel/ Error retrieving channel: ", err)
		c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, chann)
}

// @Summary      Create channel
// @Tags         Channels
// @Description  Create a new channel for a guild.
// @Accept       json
// @Produce      json
// @Param        guildID  path      string          true  "guild id"
// @Param        channel  body      models.Channel  true  "Channel values"
// @Success      201      {object}  models.Channel  "Created channel"
// @Failure      400      "Wrong values"
// @Failure      403      "Forbidden"
// @Failure      500      "Server Error"
// @Router       /guilds/{guildID}/channels [POST]
func createChannel(c echo.Context) error {
	var channel models.Channel
	guildID := c.Param("guildID")

	if err := c.Bind(&channel); err != nil || channel.GuildID != guildID {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err, "channel": channel})
	}

	_, err := db.DB.NamedExec(models.CreateChannelQuery, channel)

	if err != nil {
		// TODO switch case of sql errors.
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, channel)
}

// @Summary      Update channel values
// @Tags         Channels
// @Description  Update fields of a guild's channel
// @Accept       json
// @Produce      json
// @Param        guildID    path  string    true  "Guild id"
// @Param        channelID  path  string    true  "Channel id"
// @Success      200        "OK"  {object}  models.Channel
// @Failure      403        "Forbidden"
// @Failure      404        "Not Fountd"
// @Failure      500        "Server Error"
// @Router       /guilds/{guildID}/channels/{channelID} [PATCH]
func updateChannel(c echo.Context) error {
	guildID := c.Param("guildID")
	chanID := c.Param("id")

	var channel models.Channel

	err := db.DB.Get(&channel, "SELECT * FROM channel WHERE guild_id=? AND channel_id=?", guildID, chanID)

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

// @Summary      Delete guild channel
// @Tags         Channels
// @Description  Delete a guild channel
// @Accept       json
// @Produce      json
// @Param        guildID    path  string  true  "Guild id"
// @Param        channelID  path  string  true  "Channel id"
// @Success      206        "No Content"
// @Failure      403        "Forbidden"
// @Failure      404        "Not Fountd"
// @Failure      500        "Server Error"
// @Router       /guilds/{guildID}/channels/{channelID} [DELETE]
func deleteChannel(c echo.Context) error {
	guildID := c.Param("guildID")
	chanID := c.Param("id")

	res, err := db.DB.Exec("DELETE FROM channel WHERE guild_id = ? AND channel_id = ?", guildID, chanID)

	if err != nil {
		log.Warn("HardDeleteMember/ Error while deleting member from db: ", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not delete the channel."})
	}

	if r, _ := res.RowsAffected(); r == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}
