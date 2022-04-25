package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gyroskan/cardinal/db"
	"github.com/gyroskan/cardinal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func initUsers() {
	users := apiGroupe.Group("/users")
	users.GET("/", getUsers, isAdmin)
	users.GET("/me", getLoggedUser)
	users.GET("/:username", getUser, isAdminOrLoggedIn)
	users.PATCH("/:username", updateUser, isAdminOrLoggedIn)
	users.POST("/:username", updateAccessLvl, isAdmin)
	users.POST("/:username/ban", banUser, isAdmin)
	users.DELETE("/:username/ban", banUser, isAdmin)
	users.DELETE("/:username", deleteUser, isAdminOrLoggedIn)
}

// @Summary      Get User
// @Tags         Users
// @Description  Get a specific user
// @Param        username  path  string    true  "username"
// @Success      200       {object}  models.User "OK"
// @Failure      403       "Forbidden"
// @Failure      404       "Not found"
// @Failure      500       "Server error"
// @Router       /users/{username} [GET]
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

// @Summary      Get Logged in User
// @Tags         Users
// @Description  Get the logged in user
// @Success      200  {object}  models.User "OK"
// @Failure      403  "Forbidden"
// @Failure      500  "Server error"
// @Router       /users/me [GET]
func getLoggedUser(c echo.Context) error {
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*JwtCustomClaims)
	username := claims.Username
	var user models.User

	if err := db.DB.Get(&user, "SELECT * FROM user WHERE username=?", username); err != nil {
		log.Warn("GetLoggedUser/ Error getting user: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, user)
}

// @Summary      Get Users
// @Tags         Users
// @Description  Get a list of all existing users
// @Success      200  {array}  models.User "OK"
// @Failure      403  "Forbidden"
// @Failure      500  "Server error"
// @Router       /users/ [GET]
func getUsers(c echo.Context) error {
	var users []models.User

	if err := db.DB.Select(&users, "SELECT * FROM user"); err != nil {
		log.Warn("GetUsers/ Error getting all users: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, users)
}

// @Summary      Update User
// @Tags         Users
// @Description  Update specified User fields
// @Param        username   path  string                   true  "username"
// @Param        userModif  body  models.UserModification  true  "User modification"
// @Success      200        {object}                 models.User "OK"
// @Failure      400        "Invalid values"
// @Failure      403        "Forbidden"
// @Failure      404        "Not Found"
// @Failure      500        "Server error"
// @Router       /users/{username} [PATCH]
func updateUser(c echo.Context) error {
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

	_, err := db.DB.Exec(models.UpdateUserQuery, user)
	if err != nil {
		log.Warn("UpdateUser/ Error updating member: ", err)
		return c.JSON(http.StatusInternalServerError, "Error saving data.")
	}

	return c.JSON(http.StatusOK, user)
}

// @Summary      Update User access level
// @Tags         Users
// @Description  Update User access level
// @Param        username      path   string    true  "username"
// @Param        access_level  query  int       true  "access_level"
// @Success      200           {object}  models.User "OK"
// @Failure      403           "Forbidden"
// @Failure      404           "Not Found"
// @Failure      500           "Server error"
// @Router       /users/{username} [POST]
func updateAccessLvl(c echo.Context) error {
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

	_, err = db.DB.Exec("UPDATE user SET access_lvl=? WHERE username=?", lvl, username)
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

// @Summary      Ban User
// @Tags         Users
// @Description  Update User ban. POST to unbann, DELETE to ban.
// @Param        username  path  string  true  "username"
// @Success      200       "OK"
// @Failure      403       "Forbidden"
// @Failure      404       "Not Found"
// @Failure      500       "Server error"
// @Router       /users/{username}/ban [POST]
// @Router       /users/{username}/ban [DELETE]
func banUser(c echo.Context) error {
	username := c.Param("username")
	var tmp int
	if err := db.DB.Get(&tmp, "SELECT 1 FROM user WHERE username=?", username); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user " + username + " not found."})
		}
		log.Warn("GetUser/ Error getting user: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	ban := c.Request().Method != "POST"
	_, err := db.DB.Exec("UPDATE user SET banned=? WHERE username=?", ban, username)
	if err != nil {
		log.Warn("BanUser/error: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}

// @Summary      Delete User
// @Tags         Users
// @Description  Delete definitively the user.
// @Param        username  path  string  true  "username"
// @Success      204       "OK"
// @Failure      403       "Forbidden"
// @Failure      404       "Not Found"
// @Failure      500       "Server error"
// @Router       /users/{username} [DELETE]
func deleteUser(c echo.Context) error {
	username := c.Param("username")
	res, err := db.DB.Exec("DELETE FROM user WHERE username=?", username)
	if err != nil {
		log.Warn("deleteUser/ err: ", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if r, _ := res.RowsAffected(); r == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}
