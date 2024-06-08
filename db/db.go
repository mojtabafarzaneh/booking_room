package db

const MongoDBNameEnvName = "MONGO_DB_NAME"

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	Hotels  HotelStore
	Rooms   RoomStore
	User    UserStore
	Booking BookingStore
}
