package db

const (
	DBURI   = "mongodb://localhost:27018"
	DB_NAME = "hotel-reservation"
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

// func ToObjectID(id string) primitive.ObjectID {
// 	oid, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return oid
// }
