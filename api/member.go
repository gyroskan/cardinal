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

func initMembers() {
	g := apiGroupe.Group("/guilds/:guildID/members")
	g.GET("/", GetGuildMembers).Name = "Fetch GuildMembers."
	g.GET("/:id", GetMember).Name = "Fetch Member."
	g.POST("/", createMember).Name = "Create GuildMember."
	g.POST("/reset", resetGuildMembers).Name = "Reset Data of GuildMembers."
	g.POST("/:id/reset", resetMember).Name = "Reset Data of GuildMember."
	g.PATCH("/:id", updateMember).Name = "Update GuildMember."
	g.DELETE("/:id", hardDeleteMember).Name = "Delete GuildMember."
}

// @Summary Get Guild Members
// @Tags Members
// @Description Fetch all members of the guild.
// @Param   guildID		path	string	true	"guild id"
// @Param   limit		query	int		false	"limit to fetch" default(1)
// @Param   after		query	string	false	"higher last id fetched" default(0)
// @Success 200	"OK" {array} models.Member
// @Failure 403	"Forbidden"
// @Failure 500 "Server error"
// @Router /guilds/{guildID}/members [GET]
func GetGuildMembers(c echo.Context) error {
	guildID := c.Param("guildID")
	lastID := c.QueryParam("after")
	if lastID == "" {
		lastID = "0"
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 1
	}

	members := []models.Member{}

	err = db.DB.Select(&members, models.SelectGuildMembersQuery, guildID, lastID, limit)

	if err != nil {
		log.Warn("GetGuildMembers/ Error retrieving members from guildID: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, members)
}

// @Summary Get one Guild Member
// @Tags Members
// @Description Fetch the member of the guild.
// @Param   guildID		path	string	true	"guild id"
// @Param   memberID	path	string	true	"member id"
// @Success 200	"OK" {object} models.Member
// @Failure 403	"Forbidden"
// @Failure 404	"Not Found"
// @Failure 500 "Server error"
// @Router /guilds/{guildID}/members/{memberID} [GET]
func GetMember(c echo.Context) error {
	guildID := c.Param("guildID")
	id := c.Param("id")

	var member models.Member

	err := db.DB.Get(&member, "SELECT * FROM `member` WHERE member_id=?,guild_id=?", id, guildID)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Member with id " + id + " not found in guild " + guildID})
		}
		log.Warn("GetMember/ Error retrieving members from guildID: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, member)
}

// @Summary Create member
// @Tag Members
// @Description Create a new member from a guild.
// @Accept  json
// @Produce  json
// @Param guildID	path	string			true	"guild id"
// @Param user 		body 	models.Member 	true 	"Member values"
// @Success 201 {object} models.Member "Created member"
// @Failure 400 "Wrong values"
// @Failure 403	"Forbidden"
// @Failure 500 "Server Error"
// @Router /guilds/{guildID}/members [POST]
func createMember(c echo.Context) error {
	var member models.Member
	guildID := c.Param("guildID")

	if err := c.Bind(&member); err != nil || member.GuildID != guildID {
		return echo.NewHTTPError(http.StatusBadRequest, member)
	}

	_, err := db.DB.NamedExec(models.CreateMemberQuery, member)
	if err != nil {
		// TODO switch case of sql errors.
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, member)
}

// @Summary Reset guild's members
// @Tag Members
// @Description Reset level and xp for all guild's members.
// @Accept  json
// @Produce  json
// @Param	guildID 	path	string	true	"Guild id"
// @Success 201 "Member reset"
// @Failure 403	"Forbidden"
// @Failure 500 "Server Error"
// @Router /guilds/{guildID}/members/reset [POST]
func resetGuildMembers(c echo.Context) error {
	guildID := c.Param("guildID")

	_, err := db.DB.Exec(models.ResetGuildMembersQuery, guildID)
	if err != nil {
		log.Warn("ResetGuildMembers/ Error updating members: ", err)
		c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}

// @Summary Reset member from a guild
// @Tag Members
// @Description Reset level and xp for the specific guild's members.
// @Accept  json
// @Produce  json
// @Param	guildID 	path	string	true	"Guild id"
// @Param	memberID	path	string	true	"Guild id"
// @Success 201 "Member reset" {object} models.Member
// @Failure 403	"Forbidden"
// @Failure 500 "Server Error"
// @Router /guilds/{guildID}/members/{memberID}/reset [POST]
func resetMember(c echo.Context) error {
	guildID := c.Param("guildID")
	id := c.Param("id")

	_, err := db.DB.Exec(models.ResetGuildMembersQuery, guildID)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, "Member with id "+id+" not found in guild "+guildID)
		}
		log.Warn("ResetMember/ Error updating member: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	var memb models.Member
	err = db.DB.Get(&memb, "SELECT * FROM `member` WHERE member_id=?,guild_id=?", id, guildID)
	if err != nil {
		log.Warn("ResetMember/ Error getting member: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, memb)
}

// @Summary Update member value
// @Tag Members
// @Description Update fields of a guild's member
// @Accept  json
// @Produce  json
// @Param	guildID 	path	string	true	"Guild id"
// @Param	memberID	path	string	true	"Guild id"
// @Success 200 "OK" {object} models.Member
// @Failure 403	"Forbidden"
// @Failure 404	"Not Fountd"
// @Failure 500 "Server Error"
// @Router /guilds/{guildID}/members/{memberID} [PATCH]
func updateMember(c echo.Context) error {
	var member models.Member
	guildID := c.Param("guildID")
	id := c.Param("id")

	if err := db.DB.Get(&member, "SELECT * FROM `member` WHERE member_id=?,guild_id=?", id, guildID); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Member with id " + id + " not found in guild " + guildID})
		}
		log.Warn("UpdateMember/ Error getting member: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&member); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	member.GuildID = guildID
	member.MemberID = id

	_, err := db.DB.NamedExec(models.UpdateMemberQuery, member)
	if err != nil {
		log.Warn("UpdateMember/ Error updating member: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, member)
}

// @Summary Delete guild member
// @Tag Members
// @Description Delete a guild member
// @Accept  json
// @Produce  json
// @Param	guildID 	path	string	true	"Guild id"
// @Param	memberID	path	string	true	"Member id"
// @Success 206 "No Content"
// @Failure 403	"Forbidden"
// @Failure 404	"Not Fountd"
// @Failure 500 "Server Error"
// @Router /guilds/{guildID}/members/{memberID} [DELETE]
func hardDeleteMember(c echo.Context) error {
	guildID := c.Param("guildID")
	id := c.Param("id")

	res, err := db.DB.Exec("DELETE FROM member WHERE guild_id = ?,member_id = ?", guildID, id)

	if err != nil {
		log.Warn("HardDeleteMember/ Error while deleting member from db: ", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not delete the member."})
	}

	if r, _ := res.RowsAffected(); r == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}
