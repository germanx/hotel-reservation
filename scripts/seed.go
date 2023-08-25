package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/germanx/hotel-reservation/api"
	"github.com/germanx/hotel-reservation/db"
	"github.com/germanx/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if client.Database(db.DB_NAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.DB_NAME)
	store := &db.Store{
		User:    db.NewMongoUserStore(client, db.DB_NAME),
		Booking: db.NewMongoBookingStore(client, db.DB_NAME),
		Room:    db.NewMongoRoomStore(client, hotelStore, db.DB_NAME),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(store, "james", "foo", false)
	fmt.Println(">>> james ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println(">>> admin ->", api.CreateTokenFromUser(admin))

	hotel := fixtures.AddHotel(store, "some hotel", "bermuda", 5, nil)
	room := fixtures.AddRoom(store, "large", true, 88.44, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println(">>> booking ->", booking.ID)
}
