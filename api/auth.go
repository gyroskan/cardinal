package api

import (
	"database/sql"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/gyroskan/cardinal/db"
	"github.com/gyroskan/cardinal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	secret = os.Getenv("SECRET")
)

type (
	JwtCustomClaims struct {
		Username     string `json:"username"`
		Access_level int    `json:"access_lvl"`
		jwt.StandardClaims
	}
)

func initAuth() {
	users := apiGroupe.Group("/users")
	users.POST("/register", registerUser)
	users.POST("/login", loginUser)
}

// Register godoc
// @Summary      Register user
// @Tags         Users
// @Description  Create a new user
// @Accept       json
// @Produce      json
// @Param        user  body      models.UserCreation  true  "User values"
// @Success      201   {object}  models.User          "Created user"
// @Failure      400   "Invalid values"
// @Failure      500   "Server Error"
// @Router       /users/register [POST]
func registerUser(c echo.Context) error {
	var userCreate models.UserCreation

	if err := c.Bind(&userCreate); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := userCreate.Validate(); err != nil {
		log.Warn("Error userCreate fields: ", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	salt, err := generateSalt()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error hashing the password.")
	}

	user := models.User{
		Username:     userCreate.Username,
		Email:        userCreate.Email,
		DiscordID:    userCreate.DiscordID,
		PasswordHash: hashPassword(userCreate.Password, salt),
		Salt:         hex.EncodeToString(salt),
		AccessLvl:    2,
		CreatedAt:    time.Now(),
		Banned:       false,
	}

	if _, err := db.DB.NamedExec(models.InsertUserQuery, user); err != nil {
		switch err.(*mysql.MySQLError).Number {
		case 1062:
			return echo.NewHTTPError(http.StatusBadRequest, "Username already taken")
		default:
			log.Warn("register/ error inserting user in db: ", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary      Login user
// @Tags         Users
// @Description  Login to get user token
// @Accept       json
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        username  body      string  true  "username"
// @Param        password  body      string  true  "password"
// @Success      200       {string}  string  "Token"
// @Failure      400       "Invalid logins"
// @Failure      500       "Server Error"
// @Router       /users/login [POST]
func loginUser(c echo.Context) error {
	type userLog struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}
	var logged userLog
	if err := c.Bind(&logged); err != nil {
		log.Warn("Login/ binding error: ", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Username or password invalid")
	}
	var user models.User

	if err := db.DB.Get(&user, "SELECT * FROM user WHERE username=?", logged.Username); err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusBadRequest, "Username or password invalid")
		}
		log.Warn("GetUser/ Error getting user: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	salt, err := hex.DecodeString(user.Salt)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if user.PasswordHash != hashPassword(logged.Password, salt) {
		return echo.NewHTTPError(http.StatusBadRequest, "Username or password invalid")
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		Username:     user.Username,
		Access_level: user.AccessLvl,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 0,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Warn("Error generating token: ", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error generating token."})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
