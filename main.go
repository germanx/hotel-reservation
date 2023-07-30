package main

import (
	"context"
	"flag"
	"log"

	"github.com/germanx/hotel-reservation/api"
	"github.com/germanx/hotel-reservation/db"
	"github.com/germanx/hotel-reservation/middleware"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	// now := time.Now()
	// fmt.Println(now)
	// return

	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		// handler initialization
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
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
		apiV1 = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		auth  = app.Group("/api")
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
	apiV1.Post("/rooms", roomHandler.HandleGetRooms)
	// TODO: cancel a booking

	// bookings
	apiV1.Get("/booking", bookingHandler.HandleGetBookings)
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)

	app.Listen(*listenAddr)
}
