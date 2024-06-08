package db

const MongoDBNameEnvName = "booking-room"

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
