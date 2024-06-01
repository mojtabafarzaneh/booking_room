package db

const (
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
	DBURI      = "mongodb://192.168.1.161:27017"
)

type Store struct {
	Hotels HotelStore
	Rooms  RoomStore
	User   UserStore
}
