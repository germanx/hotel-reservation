package api

import (
	"context"
	"log"
	"testing"

	"github.com/germanx/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TEST_URI     = "mongodb://localhost:27018"
	TEST_DB_NAME = "hotel-reservation-test"
)

type testDB struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testDB) teardown(t *testing.T) {
	if err := tdb.client.Database(TEST_DB_NAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(TEST_URI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, TEST_DB_NAME)
	return &testDB{
		client: client,
		Store: &db.Store{
			Hotel:   hotelStore,
			User:    db.NewMongoUserStore(client, TEST_DB_NAME),
			Room:    db.NewMongoRoomStore(client, hotelStore, TEST_DB_NAME),
			Booking: db.NewMongoBookingStore(client, TEST_DB_NAME),
		},
	}
}
