package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/adrianowsh/booking/db"
	"github.com/adrianowsh/booking/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testdburi = "mongodb://root:1q2w3e@localhost:27017"
	dbname    = "booking_db_test"
	userColl  = "users"
)

type testdb struct {
	db.UserStoreInterface
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStoreInterface.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		t.Fatal(err)
	}
	return &testdb{
		UserStoreInterface: db.NewMongoUserStore(client, dbname),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStoreInterface)
	app.Post("/", userHandler.HandlerPostUser)

	params := types.CreateUserParams{
		FirstName: "foo",
		LastName:  "bar",
		Email:     "foobar@email.com",
		Age:       20,
		Password:  "1234567",
	}

	body, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode == fiber.StatusOK {
		var user *types.User
		json.NewDecoder(resp.Body).Decode(&user)

		if len(user.ID) == 0 {
			t.Error("expected a user id tobe set")
		}

		if len(user.Passwordhash) > 0 {
			t.Error("expected a user passwordhash noe be included in the json")
		}

		if user.FirstName != params.FirstName {
			t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
		}

		if user.LastName != params.LastName {
			t.Errorf("expected lastname %s but got %s", params.LastName, user.LastName)
		}

		if user.Email != params.Email {
			t.Errorf("expected emial %s but got %s", params.Email, user.Email)
		}

		if user.Age != params.Age {
			t.Errorf("expected age %d but got %d", params.Age, user.Age)
		}
	}
}
