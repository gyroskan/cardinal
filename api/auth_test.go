package api

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/gyroskan/cardinal/db"
	"github.com/gyroskan/cardinal/models"
	"github.com/steinfletcher/apitest"
)

var e = InitRouter()

func TestMain(m *testing.M) {
	db.Connect()
	defer db.Close()
	os.Exit(m.Run())
}

type TestDataRegister struct {
	Name string `json:"name"`
	User models.UserCreation `json:"user"`
	Expected models.User `json:"expected"`
	StatusCode int `json:"status"`
}

type TestArrayRegister struct {
	Tests []TestDataRegister `json:"tests"`
}

func TestRegister(t *testing.T) {
	sql, err := ioutil.ReadFile("../sqlScripts/builder.sql")
	if err != nil {
		t.Fatalf("Unable to read sql builder file. %v", err)
	}

	db.DB.MustExec(string(sql))

	content, err := ioutil.ReadFile("../test_data/testRegister.json")
	if err != nil {
		t.Error(err)
	}

	var tests TestArrayRegister
	err = json.Unmarshal(content, &tests)
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests.Tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Expected.CreatedAt = time.Now()
			exp, err := json.Marshal(test.Expected)
			if err != nil {
				t.Errorf("Could not Marshal expected result: %s", err)
			}

			apitest.New().
				Handler(e).
				Post(base_path + "/users/register").
				JSON(test.User).
				Expect(t).
				Status(test.StatusCode).
				Body(string(exp)).
				End()
		})
	}

	db.DB.MustExec(string(sql))
}
//
//func TestRegisterOld(t *testing.T) {
//	user := models.UserCreation{
//		Username: "test1",
//		Email:    "test@mail.com",
//		Password: "123",
//	}
//
//	expected := models.User{
//		Username:  user.Username,
//		Email:     user.Email,
//		DiscordID: user.DiscordID,
//		AccessLvl: 2,
//		CreatedAt: time.Now(),
//		Banned:    false,
//	}
//
//	data, err := json.Marshal(user)
//	if err != nil {
//		t.Skip("Failed to marshal user")
//	}
//	d, err := json.Marshal(expected)
//	if err != nil {
//		t.Skip("Failed to marshal user")
//	}
//
//	// valid, discordID null
//	apitest.New().
//		Handler(e).
//		Post(base_path + "/users/register").
//		JSON(data).
//		Expect(t).
//		Status(http.StatusOK).
//		Body(string(d)).
//		End()
//	
//	apitest.New().
//		Handler(e).
//		Post(base_path + "/users/register").
//		JSON(data).
//		Expect(t).
//		Status(http.StatusOK).
//		Body(string(d)).
//		End()
//}
//