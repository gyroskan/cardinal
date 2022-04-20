package api

import (
	"database/sql"
	"net/http"

	"github.com/gyroskan/cardinal/db"
	"github.com/gyroskan/cardinal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func initBans() {
	b := apiGroupe.Group("/guilds/:guildID/members/:memberID/bans")
	b.GET("/", getBans).Name = "Fetch all bans of a member."
	b.GET("/:banID", getBan).Name = "Fetch a ban of a member."
	b.POST("/", createBan).Name = "Create a ban for a member."
	b.DELETE("/:banID", deleteBan).Name = "Delete a ban of a member."
}

// @Summary      Get Member Bans
// @Tags         Bans
// @Description  Fetch all bans of the member.
// @Param        guildID   path  string   true  "guild id"
// @Param        memberID  path  string   true  "member id"
// @Success      200       "OK"  {array}  models.Ban
// @Failure      403       "Forbidden"
// @Failure      500       "Server error"
// @Router       /guilds/{guildID}/members/{memberID}/bans [GET]
func getBans(c echo.Context) error {
	guildID := c.Param("guildID")
	memberID := c.Param("guildID")
	var bans []models.Ban

	err := db.DB.Select(&bans, "SELECT * FROM ban WHERE guild_id=? AND member_id=?", guildID, memberID)

	if err != nil {
		log.Warn("GetBans/ Error retrieving bans: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, bans)
}

// @Summary      Get one ban
// @Tags         Bans
// @Description  Fetch the ban of the member.
// @Param        guildID   path  string    true  "guild id"
// @Param        memberID  path  string    true  "member id"
// @Param        banID     path  string    true  "ban id"
// @Success      200       "OK"  {object}  models.Ban
// @Failure      403       "Forbidden"
// @Failure      404       "Not Found"
// @Failure      500       "Server error"
// @Router       /guilds/{guildID}/members/{memberID}/bans/{banID} [GET]
func getBan(c echo.Context) error {
	guildID := c.Param("guildID")
	memberID := c.Param("guildID")
	banID := c.Param("banID")
	var ban models.Ban

	err := db.DB.Get(&ban, "SELECT * FROM ban WHERE guild_id=? AND member_id=? AND ban_id=?", guildID, memberID, banID)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, nil)
		}
		log.Warn("GetBan/ Error retrieving Ban: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, ban)
}

// @Summary      Create ban
// @Tags         Bans
// @Description  Create a new ban for a member.
// @Accept       json
// @Produce      json
// @Param        guildID   path      string      true  "guild id"
// @Param        memberID  path      string      true  "member id"
// @Success      201       {object}  models.Ban  "Created role"
// @Failure      400       "Wrong values"
// @Failure      403       "Forbidden"
// @Failure      500       "Server Error"
// @Router       /guilds/{guildID}/members/{memberID}/bans [POST]
func createBan(c echo.Context) error {
	guildID := c.Param("guildID")
	memberID := c.Param("guildID")
	var ban models.Ban

	if err := c.Bind(&ban); err != nil || ban.GuildID != guildID || ban.MemberID != memberID {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err, "ban": ban})
	}

	query, err := db.DB.PrepareNamed(models.CreateBanQuery)

	if err != nil {
		log.Error("CreateBan/ error while preparing query:", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	res, err := query.Exec(ban)

	if err != nil {
		log.Error("CreateBan/ error while executing query:", err)
		// TODO switch case of sql errors.
		return c.JSON(http.StatusInternalServerError, nil)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Error("CreateBan/ Error while getting last index: ", err)
		// TODO switch case of sql errors.
		return c.JSON(http.StatusInternalServerError, nil)
	}

	ban.BanID = int(id)
	return c.JSON(http.StatusCreated, ban)
}

// @Summary      Delete member's ban
// @Tags         Bans
// @Description  Delete a member's ban
// @Accept       json
// @Produce      json
// @Param        guildID   path  string  true  "Guild id"
// @Param        memberID  path  string  true  "member id"
// @Param        banID     path  string  true  "ban id"
// @Success      206       "No Content"
// @Failure      403       "Forbidden"
// @Failure      404       "Not Found"
// @Failure      500       "Server Error"
// @Router       /guilds/{guildID}/members/{memberID}/bans/{banID} [DELETE]
func deleteBan(c echo.Context) error {
	guildID := c.Param("guildID")
	memberID := c.Param("guildID")
	banID := c.Param("banID")

	res, err := db.DB.Exec("DELETE FROM ban WHERE guild_id = ? AND member_id= ? AND ban_id=?", guildID, memberID, banID)

	if err != nil {
		log.Error("DeleteBan/ Error while deleting ban from db: ", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not delete the ban."})
	}

	if r, _ := res.RowsAffected(); r == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}
