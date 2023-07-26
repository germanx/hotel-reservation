package db

const (
	DBURI      = "mongodb://localhost:27018"
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}

// func ToObjectID(id string) primitive.ObjectID {
// 	oid, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return oid
// }
