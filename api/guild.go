package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gyroskan/cardinal/db"
	"github.com/gyroskan/cardinal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func initGuildGroup() {
	g := apiGroupe.Group("/guild")
	g.GET("/", getGuilds).Name = "Fetch All Guilds."
	g.GET("/:id", getGuild).Name = "Fetch Guild by id."
	g.POST("/", createGuild).Name = "Create new guild."
	// g.PATCH("/:id", updateGuild).Name = "Update guild."
	// g.POST("/:id/reset", resetGuild).Name = "Reset guild."
	// g.DELETE("/:id", hardDeleteGuild).Name = "Hard Delete guild."
}

// @Summary Get Guilds
// @Tags Guilds
// @Description Fetch all guilds.
// @Success 200	"OK" {array} models.Guild
// @Failure 403	"Forbidden"
// @Failure 500 "Server error"
// @Router /guilds [GET]
func getGuilds(c echo.Context) error {
	var guilds []models.Guild

	err := db.DB.Select(&guilds, "SELECT * FROM `guild`")

	if err != nil {
		log.Warn("GetGuilds/ Error retrieving guilds: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, guilds)
}

// @Summary Get guild
// @Tags Guilds
// @Description Fetch a specific guild
// @Param   guildID	path	string	true	"guild id"
// @Param 	members	query	bool	false	"fetch members"
// @Success 200	"OK" {object} models.Guild
// @Failure 403	"Forbidden"
// @Failure 404	"Not Found"
// @Failure 500 "Server error"
// @Router /guilds/{guildID} [GET]
func getGuild(c echo.Context) error {
	id := c.Param("id")
	members, err := strconv.ParseBool(c.QueryParam("members"))

	if err != nil {
		members = false
	}

	var guild models.Guild
	err = db.DB.Get(&guild, "SELECT * FROM guild WHERE `guild_id`=?", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, nil)
		}

		log.Warn("GetGuild/ error retrieving guild: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if members {
		err = db.DB.Select(&guild.Members, "SELECT * FROM member WHERE `guild_id`=?", id)

		if err != nil {
			log.Warn("GetGuild/ Error retrieving members of guild: ", err)
			guild.Members = nil
		}
	}

	return c.JSON(http.StatusOK, guild)
}

// @Summary Create guild
// @Tags Guilds
// @Description Create a new Guild
// @Success 201	"Created" {object} models.Guild
// @Failure 400	"Bad Request"
// @Failure 403	"Forbidden"
// @Failure 409	"Conflict"
// @Failure 500 "Server error"
// @Router /guilds/ [POST]
func createGuild(c echo.Context) error {
	guild := new(models.Guild)
	if err := c.Bind(&guild); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, guild)
	}

	_, err := db.DB.NamedExec(models.CreateGuildQuery, guild)

	if err != nil {
		if err == nil { // TODO error code for conflict
			return c.JSON(http.StatusConflict, echo.Map{"message": "The guild with id " + guild.GuildID + " already exists."})
		}
		log.Warn("CreateGuild/ Error inserting values: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	for _, memb := range guild.Members {
		if memb.GuildID == guild.GuildID {
			_, err = db.DB.NamedExec(models.CreateMemberQuery, memb)
			if err != nil {
				//TODO ignore existing member
				log.Error("CreateGuild/ Error while creating member in DB: ", err)
			}
		}
	}

	return c.JSON(http.StatusCreated, guild)
}
