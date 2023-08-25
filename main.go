package main

import (
	"context"
	"flag"
	"log"

	"github.com/germanx/hotel-reservation/api"
	"github.com/germanx/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		// handler initialization
		hotelStore   = db.NewMongoHotelStore(client, db.DB_NAME)
		roomStore    = db.NewMongoRoomStore(client, hotelStore, db.DB_NAME)
		userStore    = db.NewMongoUserStore(client, db.DB_NAME)
		bookingStore = db.NewMongoBookingStore(client, db.DB_NAME)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)

		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiV1 = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin = apiV1.Group("/admin", api.AdminAuth)
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// user handlers
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Get("/users", userHandler.HandleGetUsers)

	// hotel handlers
	apiV1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// rooms
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiV1.Get("/rooms", roomHandler.HandleGetRooms)
	// TODO: cancel a booking

	// bookings
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiV1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBookings)

	// admin handlers
	admin.Get("/bookings", bookingHandler.HandleGetBookings)

	app.Listen(*listenAddr)
}
