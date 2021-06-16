package api

import (
	"database/sql"
	"net/http"

	"github.com/gyroskan/cardinal/db"
	"github.com/gyroskan/cardinal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func initWarn() {
	w := apiGroupe.Group("/guilds/:guildID/members/:memberID/warns")
	w.GET("/", getWarns)
	w.GET("/:warnID", getWarn)
	w.POST("/", createWarn)
	w.DELETE("/:warnID", deleteWarn)
}

// @Summary Get Member Warns
// @Tags Warns
// @Description Fetch all warns of the member.
// @Param   guildID		path	string	true	"guild id"
// @Param   memberID	path	string	true	"member id"
// @Success 200	"OK" {array} models.Warn
// @Failure 403	"Forbidden"
// @Failure 500 "Server error"
// @Router /guilds/{guildID}/members/{memberID}/warns [GET]
func getWarns(c echo.Context) error {
	guildID := c.Param("guildID")
	memberID := c.Param("guildID")
	var bans []models.Warn

	err := db.DB.Select(&bans, "SELECT * FROM warn WHERE guild_id=? AND member_id=?", guildID, memberID)

	if err != nil {
		log.Warn("GetWarns/ Error retrieving warns: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, bans)
}

// @Summary Get one warn
// @Tags Warns
// @Description Fetch the warn of the member.
// @Param   guildID		path	string	true	"guild id"
// @Param   memberID	path	string	true	"member id"
// @Param   warnID		path	string	true	"warn id"
// @Success 200	"OK" {object} models.Warn
// @Failure 403	"Forbidden"
// @Failure 404	"Not Found"
// @Failure 500 "Server error"
// @Router /guilds/{guildID}/members/{memberID}/warns/{warnID} [GET]
func getWarn(c echo.Context) error {
	guildID := c.Param("guildID")
	memberID := c.Param("guildID")
	warnID := c.Param("warnID")
	var warn models.Warn

	err := db.DB.Get(&warn, "SELECT * FROM warn WHERE guild_id=? AND member_id=? AND warn_id=?", guildID, memberID, warnID)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, nil)
		}
		log.Warn("Getwarn/ Error retrieving warn: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, warn)
}

// @Summary Create warn
// @Tag Warns
// @Description Create a new warn for a member.
// @Accept  json
// @Produce  json
// @Param   guildID		path	string	true	"guild id"
// @Param   memberID	path	string	true	"member id"
// @Success 201 {object} models.Warn "Created warn"
// @Failure 400 "Wrong values"
// @Failure 403	"Forbidden"
// @Failure 500 "Server Error"
// @Router /guilds/{guildID}/members/{memberID}/warns [POST]
func createWarn(c echo.Context) error {
	guildID := c.Param("guildID")
	memberID := c.Param("guildID")
	var warn models.Warn

	if err := c.Bind(&warn); err != nil || warn.GuildID != guildID || warn.MemberID != memberID {
		return echo.NewHTTPError(http.StatusBadRequest, warn)
	}

	query, err := db.DB.PrepareNamed(models.CreateWarnQuery)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	err = query.Get(&warn.WarnID, warn)

	if err != nil {
		// TODO switch case of sql errors.
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, warn)
}

// @Summary Delete member's warn
// @Tag Warns
// @Description Delete a member's warn
// @Accept  json
// @Produce  json
// @Param	guildID 	path	string	true	"Guild id"
// @Param   memberID	path	string	true	"member id"
// @Param   warnID		path	string	true	"warn id"
// @Success 206 "No Content"
// @Failure 403	"Forbidden"
// @Failure 404	"Not Found"
// @Failure 500 "Server Error"
// @Router /guilds/{guildID}/members/{memberID}/warns/{warnID} [DELETE]
func deleteWarn(c echo.Context) error {
	guildID := c.Param("guildID")
	memberID := c.Param("guildID")
	warnID := c.Param("warnID")

	res, err := db.DB.Exec("DELETE FROM warn WHERE guild_id=? AND member_id=? AND warn_id=?", guildID, memberID, warnID)

	if err != nil {
		log.Error("DeleteWarn/ Error while deleting warn from db: ", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not delete the warn."})
	}

	if r, _ := res.RowsAffected(); r == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}
