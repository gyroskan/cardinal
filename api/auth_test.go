package api

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gyroskan/cardinal/models"
	"github.com/steinfletcher/apitest"
)

var e = InitRouter()

type TestDataRegister struct {
	
}

func TestRegister(t *testing.T) {

}

func TestRegisterOld(t *testing.T) {
	user := models.UserCreation{
		Username: "test1",
		Email:    "test@mail.com",
		Password: "123",
	}

	expected := models.User{
		Username:  user.Username,
		Email:     user.Email,
		DiscordID: user.DiscordID,
		AccessLvl: 2,
		CreatedAt: time.Now(),
		Banned:    false,
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Skip("Failed to marshal user")
	}
	d, err := json.Marshal(expected)
	if err != nil {
		t.Skip("Failed to marshal user")
	}

	// valid, discordID null
	apitest.New().
		Handler(e).
		Post(base_path + "/users/register").
		JSON(data).
		Expect(t).
		Status(http.StatusOK).
		Body(string(d)).
		End()
	
	apitest.New().
		Handler(e).
		Post(base_path + "/users/register").
		JSON(data).
		Expect(t).
		Status(http.StatusOK).
		Body(string(d)).
		End()
}
