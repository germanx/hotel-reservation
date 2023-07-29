package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/germanx/hotel-reservation/db"
	"github.com/germanx/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TEST_URI = "mongodb://localhost:27018"
)

type testDB struct {
	db.UserStore
}

func (tdb *testDB) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(TEST_URI))
	if err != nil {
		log.Fatal(err)
	}
	return &testDB{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestPostUser(t *testing.T) {
	// t.Fail()
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "some@foo.com",
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Password:  "1234567",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	// fmt.Println(user)
	if len(user.ID) == 0 {
		t.Error("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Error("expecting a password not to be included in response")
	}

	if user.FirstName != params.FirstName {
		t.Errorf("expected FirstName %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected LastName %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected Email %s but got %s", params.Email, user.Email)
	}

	// bb, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(bb))
	// fmt.Println(resp.Status)
}
