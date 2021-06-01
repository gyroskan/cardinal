package api

import (
	"database/sql"
	"net/http"

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
