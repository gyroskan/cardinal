package api

import (
	"database/sql"
	"net/http"

	"github.com/gyroskan/cardinal/db"
	"github.com/gyroskan/cardinal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func initRoles() {
	r := apiGroupe.Group("/guilds/:guildID/roles")
	r.GET("/", getRoles).Name = "Fetch all guild roles."
	r.GET("/:id", getRole).Name = "Fetch a guild role."
	r.POST("/", createRole).Name = "Create a guild role."
}

// @Summary Get Guild roles
// @Tags Roles
// @Description Fetch all roles of the guild.
// @Param   guildID		path	string	true	"guild id"
// @Success 200	"OK" {array} models.Roles
// @Failure 403	"Forbidden"
// @Failure 500 "Server error"
// @Router /guilds/{guildID}/roles [GET]
func getRoles(c echo.Context) error {
	guildID := c.Param("guildID")
	var roles []models.Role

	err := db.DB.Select(&roles, "SELECT * FROM role WHERE guild_id=?", guildID)

	if err != nil {
		log.Warn("GetRoles/ Error retrieving roles: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, roles)
}

// @Summary Get one Guild role
// @Tags Roles
// @Description Fetch the role of the guild.
// @Param   guildID		path	string	true	"guild id"
// @Param   roleID	path	string	true	"role id"
// @Success 200	"OK" {object} models.Role
// @Failure 403	"Forbidden"
// @Failure 404	"Not Found"
// @Failure 500 "Server error"
// @Router /guilds/{guildID}/roles/{roleID} [GET]
func getRole(c echo.Context) error {
	guildID := c.Param("guildID")
	roleID := c.Param("id")
	var role models.Role

	err := db.DB.Get(&role, "SELECT * FROM role WHERE guild_id=?,role_id", guildID, roleID)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, nil)
		}
		log.Warn("GetRoles/ Error retrieving roles: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, role)
}

// @Summary Create role
// @Tag Roles
// @Description Create a new role for a guild.
// @Accept  json
// @Produce  json
// @Param   guildID	path	string			true	"guild id"
// @Param 	role	body 	models.Role 	true 	"Role values"
// @Success 201 {object} models.Role "Created role"
// @Failure 400 "Wrong values"
// @Failure 403	"Forbidden"
// @Failure 500 "Server Error"
// @Router /guilds/{guildID}/roles [POST]
func createRole(c echo.Context) error {
	var role models.Role
	guildID := c.Param("guildID")

	if err := c.Bind(&role); err != nil || role.GuildID != guildID {
		return echo.NewHTTPError(http.StatusBadRequest, role)
	}

	_, err := db.DB.NamedExec(models.CreateRoleQuery, role)

	if err != nil {
		// TODO switch case of sql errors.
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, role)
}
