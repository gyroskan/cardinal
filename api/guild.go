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

func initGuilds() {
	g := apiGroupe.Group("/guilds")
	g.GET("/", getGuilds).Name = "Fetch All Guilds."
	g.GET("/:id", getGuild).Name = "Fetch Guild by id."
	g.POST("/", createGuild).Name = "Create new guild."
	g.PATCH("/:id", updateGuild).Name = "Update guild."
	g.POST("/:id/reset", resetGuild).Name = "Reset guild."
	g.DELETE("/:id", hardDeleteGuild).Name = "Hard Delete guild."
}

// @Summary      Get all Guilds
// @Tags         Guilds
// @Description  Fetch all guilds.
// @Success      200  {array}  models.Guild  "OK"
// @Failure      403  "Forbidden"
// @Failure      500  "Server error"
// @Router       /guilds/ [GET]
func getGuilds(c echo.Context) error {
	var guilds []models.Guild

	err := db.DB.Select(&guilds, "SELECT * FROM `guild`")

	if err != nil {
		log.Warn("GetGuilds/ Error retrieving guilds: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, guilds)
}

// @Summary      Get one guild
// @Tags         Guilds
// @Description  Fetch a specific guild
// @Param        guildID  path      string        true  "guild id"
// @Success      200      {object}  models.Guild  "OK"
// @Failure      403      "Forbidden"
// @Failure      404      "Not Found"
// @Failure      500      "Server error"
// @Router       /guilds/{guildID} [GET]
func getGuild(c echo.Context) error {
	id := c.Param("id")
	// members, err := strconv.ParseBool(c.QueryParam("members"))

	// if err != nil {
	// 	members = false
	// }

	var guild models.Guild
	err := db.DB.Get(&guild, "SELECT * FROM guild WHERE `guild_id`=?", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, nil)
		}

		log.Warn("GetGuild/ error retrieving guild: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// if members {
	// 	err = db.DB.Select(&guild.Members, "SELECT * FROM member WHERE `guild_id`=?", id)

	// 	if err != nil {
	// 		log.Warn("GetGuild/ Error retrieving members of guild: ", err)
	// 		guild.Members = nil
	// 	}
	// }

	return c.JSON(http.StatusOK, guild)
}

// @Summary      Create guild
// @Tags         Guilds
// @Description  Creates a new Guild with the given values
// @param        guild  body      models.Guild  true  "Provide the guild values"
// @Success      201    {object}  models.Guild  "Created"
// @Failure      400    "Bad Request"
// @Failure      403    "Forbidden"
// @Failure      409    "Conflict"
// @Failure      500    "Server error"
// @Router       /guilds/ [POST]
func createGuild(c echo.Context) error {
	var guild models.Guild
	if err := c.Bind(&guild); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err, "guild": guild})
	}

	_, err := db.DB.NamedExec(models.CreateGuildQuery, guild)

	if err != nil {
		if err == nil { // TODO error code for conflict
			return c.JSON(http.StatusConflict, echo.Map{"message": "The guild with id " + guild.GuildID + " already exists."})
		}
		log.Warn("CreateGuild/ Error inserting values: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, guild)
}

// @Summary      Update guild values
// @Tags         Guilds
// @Description  Update fields of a guild
// @Accept       json
// @Produce      json
// @Param        guildID  path      string        true  "Guild id"
// @Param        guild    body      models.Guild  true  "Guild modifications"
// @Success      200      {object}  models.Guild  "OK"
// @Failure      403      "Forbidden"
// @Failure      404      "Not Found"
// @Failure      500      "Server Error"
// @Router       /guilds/{guildID} [PATCH]
func updateGuild(c echo.Context) error {
	id := c.Param("id")
	var guild models.Guild

	if err := db.DB.Get(&guild, "SELECT * FROM guild WHERE `guild_id`=?", id); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Guild with id`" + id + "` not found."})
		}
		log.Warn("UpdateGuild/ Error while retrieving guild: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// If some fields were not provided, the previous value are kept.
	if err := json.NewDecoder(c.Request().Body).Decode(&guild); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	guild.GuildID = id

	_, err := db.DB.NamedExec(models.UpdateGuildQuery, guild)

	if err != nil {
		log.Warn("UpdateGuild/ Error while Updating DB: ", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, guild)
}

// @Summary      Reset guild
// @Tags         Guilds
// @Description  Reset guild parameters to default values.
// @Description  Do not change members values.
// @Accept       json
// @Produce      json
// @Param        guildID  path      string        true  "Guild id"
// @Success      200      {object}  models.Guild  "OK"
// @Failure      403      "Forbidden"
// @Failure      500      "Server Error"
// @Router       /guilds/{guildID}/reset [POST]
func resetGuild(c echo.Context) error {
	guildID := c.Param("guildID")

	_, err := db.DB.Exec(models.ResetGuildQuery, guildID)

	if err != nil {
		if err == sql.ErrNoRows { //should not happened
			return c.JSON(http.StatusNotFound, "Guild with id "+guildID+" not found.")
		}
		log.Error("ResetGuild/ Error updating guild: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	var guild models.Guild
	err = db.DB.Get(&guild, "SELECT * FROM guild WHERE `guild_id`=?", guildID)

	if err != nil {
		log.Error("GetGuild/ error retrieving guild: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, guild)
}

// @Summary      Delete guild
// @Tags         Guilds
// @Description  Delete a guild
// @Accept       json
// @Produce      json
// @Param        guildID  path  string  true  "Guild id"
// @Success      204      "No Content"
// @Failure      403      "Forbidden"
// @Failure      404      "Not Found"
// @Failure      500      "Server Error"
// @Router       /guilds/{guildID} [DELETE]
func hardDeleteGuild(c echo.Context) error {
	id := c.Param("id")

	res, err := db.DB.Exec("DELETE FROM guild WHERE guild_id = ?", id)

	if err != nil {
		log.Error("HardDeleteGuild/ Error while deleting guild from db: ", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not delete the guild."})
	}

	if r, _ := res.RowsAffected(); r == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}
