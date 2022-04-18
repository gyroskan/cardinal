package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	r.PATCH("/:id", updateRole).Name = "Update a guild role."
	r.DELETE("/:id", deleteRole).Name = "Delete a guild role."
}

// @Summary      Get Guild roles
// @Tags         Roles
// @Description  Fetch all roles of the guild.
// @Param        guildID             path     string  true                      "guild id"
// @Param        reward              query    bool    false                     "reward from this lvl only"  default(false)
// @Param        ignored             query    bool    false                     "ignored roles only"             default(false)
// @Param        xpBlacklist  query  bool     false   "xpBlacklist roles only"  default(false)
// @Success      200          "OK"   {array}  models.Roles
// @Failure      403          "Forbidden"
// @Failure      500          "Server error"
// @Router       /guilds/{guildID}/roles [GET]
func getRoles(c echo.Context) error {
	guildID := c.Param("guildID")
	ignored := false
	xpBlacklisted := false
	reward := 0
	if c.QueryParam("ignored") != "" {
		ignored, _ = strconv.ParseBool(c.QueryParam("ignored"))
	}
	if c.QueryParam("xpBlacklist") != "" {
		xpBlacklisted, _ = strconv.ParseBool(c.QueryParam("xpBlacklist"))
	}
	if c.QueryParam("reward") != "" {
		reward, _ = strconv.Atoi(c.QueryParam("reward"))
	}
	var roles []models.Role

	query := "SELECT * FROM role WHERE guild_id=?"
	if ignored {
		query += " AND ignored=true"
	}
	if xpBlacklisted {
		query += " AND xp_blacklisted=true"
	}
	if reward != 0 {
		query += fmt.Sprintf(" AND reward=%d", reward)
	}

	err := db.DB.Select(&roles, query, guildID)

	if err != nil {
		log.Warn("GetRoles/ Error retrieving roles: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, roles)
}

// @Summary      Get one Guild role
// @Tags         Roles
// @Description  Fetch the role of the guild.
// @Param        guildID        path      string  true  "guild id"
// @Param        roleID   path  string    true    "role id"
// @Success      200      "OK"  {object}  models.Role
// @Failure      403      "Forbidden"
// @Failure      404      "Not Found"
// @Failure      500      "Server error"
// @Router       /guilds/{guildID}/roles/{roleID} [GET]
func getRole(c echo.Context) error {
	guildID := c.Param("guildID")
	roleID := c.Param("id")
	var role models.Role

	err := db.DB.Get(&role, "SELECT * FROM role WHERE guild_id=? AND role_id=?", guildID, roleID)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, nil)
		}
		log.Warn("GetRoles/ Error retrieving roles: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, role)
}

// @Summary      Create role
// @Tag          Roles
// @Description  Create a new role for a guild.
// @Accept       json
// @Produce      json
// @Param        guildID  path      string                      true  "guild id"
// @Param                 role      body           models.Role        true    "Role values"
// @Success      201      {object}  models.Role  "Created role"
// @Failure      400      "Wrong values"
// @Failure      403      "Forbidden"
// @Failure      500      "Server Error"
// @Router       /guilds/{guildID}/roles [POST]
func createRole(c echo.Context) error {
	var role models.Role
	guildID := c.Param("guildID")

	if err := c.Bind(&role); err != nil || role.GuildID != guildID {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err, "role": role})
	}

	_, err := db.DB.NamedExec(models.CreateRoleQuery, role)

	if err != nil {
		// TODO switch case of sql errors.
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, role)
}

// @Summary      Update role values
// @Tag          Roles
// @Description  Update fields of a guild's role
// @Accept       json
// @Produce      json
// @Param        guildID        path      string  true  "Guild id"
// @Param        roleID   path  string    true    "role id"
// @Success      200      "OK"  {object}  models.Role
// @Failure      403      "Forbidden"
// @Failure      404      "Not Fountd"
// @Failure      500      "Server Error"
// @Router       /guilds/{guildID}/roles/{roleID} [PATCH]
func updateRole(c echo.Context) error {
	guildID := c.Param("guildID")
	roleID := c.Param("id")
	var role models.Role

	err := db.DB.Get(&role, "SELECT * FROM role WHERE guild_id=? AND role_id=?", guildID, roleID)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, nil)
		}
		log.Warn("UpdateRole/ Error retrieving roles: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&role); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	role.GuildID = guildID
	role.RoleID = roleID

	_, err = db.DB.NamedExec(models.UpdateRoleQuery, role)
	if err != nil {
		log.Error("UpdateRole/ Error updating role: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, role)
}

// @Summary      Delete guild role
// @Tag          Roles
// @Description  Delete a guild role
// @Accept       json
// @Produce      json
// @Param        guildID  path  string  true  "Guild id"
// @Param        roleID   path  string  true  "role id"
// @Success      206      "No Content"
// @Failure      403      "Forbidden"
// @Failure      404      "Not Fountd"
// @Failure      500      "Server Error"
// @Router       /guilds/{guildID}/roles/{roleID} [DELETE]
func deleteRole(c echo.Context) error {
	guildID := c.Param("guildID")
	roleID := c.Param("id")

	res, err := db.DB.Exec("DELETE FROM role WHERE guild_id = ? AND role_id = ?", guildID, roleID)

	if err != nil {
		log.Error("HardDeleteRole/ Error while deleting role from db: ", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not delete the role."})
	}

	if r, _ := res.RowsAffected(); r == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}
