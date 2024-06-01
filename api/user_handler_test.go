package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/db"
	"github.com/mojtabafarzaneh/hotel_reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	store db.Store
}

func (tdb *testdb) teardrop(t *testing.T) {
	if err := tdb.store.User.Drop(context.TODO()); err != nil {
		t.Fatal()
	}

}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		store: db.Store{
			User: db.NewMongoUserStore(client),
		},
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardrop(t)

	app := fiber.New()
	UserHandler := NewUserHandler(tdb.store)
	app.Post("/", UserHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "moji",
		LastName:  "farzaneh",
		Email:     "email@emial.com",
		Password:  "faafuhgawe",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(res.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("expected user ID but got none")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected user firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("excepted user lastname %s but got %s", params.LastName, user.LastName)

	}
	if user.Email != params.Email {
		t.Errorf("excepted user email %s but got %s", params.Email, user.Email)
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting encrypted password to not be included in the response")
	}
}
