package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/db"
	"github.com/mojtabafarzaneh/hotel_reservation/types"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "moji",
		LastName:  "lname",
		Email:     "moj@gmail.com",
		Password:  "supersecurepass",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = userStore.InsertUsers(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthenticationSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardrop(t)
	insertedUser := insertTestUser(t, tdb.Store.User)

	app := fiber.New()
	AuthHandler := NewAuthHnadler(tdb.Store.User)
	app.Post("/", AuthHandler.HandleAuthentication)

	params := AuthParams{
		Email:    "moj@gmail.com",
		Password: "supersecurepass",
	}
	p, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/", bytes.NewReader(p))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("expected status 200 but got", res.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(res.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if len(authResp.Token) == 0 {
		t.Fatal("expected jwt token but got none")
	}
	//set the password to an empty string cause we do not back passwords to the client!
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatal("expected the inserted user but got the wrong one")
	}
}
