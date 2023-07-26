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

// const userColl = "users"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
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
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		userStore  = db.NewMongoUserStore(client)
		store      = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}
		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		authHandler  = api.NewAuthHandler(userStore)
		app          = fiber.New(config)
		apiV1        = app.Group("/api/v1", middleware.JWTAuthentication)
		auth         = app.Group("/api")
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

	app.Listen(*listenAddr)
}

// func addUser() {
// 	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ctx := context.Background()
// 	coll := client.Database(dbName).Collection(userColl)

// 	user := types.User{
// 		FirstName: "James",
// 		LastName:  "Watercooler",
// 	}
// 	_, err = coll.InsertOne(ctx, user)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var james types.User
// 	if err := coll.FindOne(ctx, bson.M{}).Decode(&james); err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(james)
// }
