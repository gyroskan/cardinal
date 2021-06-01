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

func initUsers() {
	users := apiGroupe.Group("/users")
	users.GET("/", getUsers, isAdmin)
	users.GET("/:username", getUser, isAdminOrLoggedIn)
}

// @Summary Get User
// @Tags Users
// @Description Get a specific user
// @Param   username	path	string	true	"username"
// @Success 200	"OK" {object} models.User
// @Failure 403	"Forbidden"
// @Failure 404	"Not found"
// @Failure 500 "Server error"
// @Router /users/{username} [GET]
func getUser(c echo.Context) error {
	username := c.Param("username")
	var user models.User

	if err := db.DB.Get(&user, "SELECT * FROM user WHERE username=?", username); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user " + username + " not found."})
		}
		log.Warn("GetUser/ Error getting user: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, user)
}

// @Summary Get Users
// @Tags Users
// @Description Get a list of existing users
// @Success 200	"OK" {array} models.User
// @Failure 403	"Forbidden"
// @Failure 500 "Server error"
// @Router /users/ [GET]
func getUsers(c echo.Context) error {
	var users []models.User

	if err := db.DB.Select(&users, "SELECT * FROM user"); err != nil {
		log.Warn("GetUsers/ Error getting all users: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, users)
}

// @Summary Update User
// @Tags Users
// @Description Update specified User fields
// @Param   username		path	string	true	"username"
// @Param   userModif		body	{object} models.UserModification	true	"fields to modify"
// @Success 200	"OK" {object} models.User
// @Failure 400 "Invalid values"
// @Failure 403	"Forbidden"
// @Failure 404	"NotFound"
// @Failure 500 "Server error"
// @Router /users/{username} [PATCH]
func UpdateUser(c echo.Context) error {
	var user models.User
	username := c.Param("username")

	if err := db.DB.Get(&user, "SELECT * FROM `user` WHERE username=?", username); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user " + username + " not found."})
		}
		log.Warn("UpdateUser/ Error getting user: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	userUp := models.UserModification{
		Email:     user.Email,
		DiscordID: user.DiscordID,
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&userUp); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if userUp.OldPassword != "" {
		if hashPassword(userUp.OldPassword, []byte(user.Salt)) != user.PasswordHash {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid old password.")
		}
		salt, err := generateSalt()
		if err != nil {
			log.Warn("UpdateUser/ Error generating salt: ", err)
			return c.JSON(http.StatusInternalServerError, "Error hashing new password.")
		}
		user.PasswordHash = hashPassword(userUp.Password, salt)
		user.Salt = string(salt)
	}

	_, err := db.DB.Query(models.UpdateUserQuery, user)
	if err != nil {
		log.Warn("UpdateUser/ Error updating member: ", err)
		return c.JSON(http.StatusInternalServerError, "Error saving data.")
	}

	return c.JSON(http.StatusOK, user)
}

// @Summary Update User access level
// @Tags Users
// @Description Update User access level
// @Param   username		path	string	true	"username"
// @Param   access_level	query	int		true	"access_level"
// @Success 200	"OK" {object} models.User
// @Failure 403	"Forbidden"
// @Failure 404	"NotFound"
// @Failure 500 "Server error"
// @Router /users/{username} [POST]
func UpdateAccessLvl(c echo.Context) error {
	lvl, err := strconv.Atoi(c.QueryParam("access_level"))
	if err != nil || lvl > 2 || lvl < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid access_level")
	}

	username := c.Param("username")
	var user models.User
	if err := db.DB.Get(&user, "SELECT * FROM `user` WHERE username=?", username); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user " + username + " not found."})
		}
		log.Warn("UpdateAccessLvl/ Error getting user: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	_, err = db.DB.Queryx("UPDATE user SET access_lvl=? WHERE username=?", lvl, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user " + username + " not found."})
		}
		log.Warn("UpdateAccessLvl/ Error updating member: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error saving data.")
	}

	user.AccessLvl = lvl

	return c.JSON(http.StatusOK, user)
}
