package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mojtabafarzaneh/hotel_reservation/api/middleware"
	"github.com/mojtabafarzaneh/hotel_reservation/db/fixtures"
	"github.com/mojtabafarzaneh/hotel_reservation/types"
)

func TestGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardrop(t)

	var (
		noneAuthUser   = fixtures.AddUser(db.Store, "ri", "stu", false)
		user           = fixtures.AddUser(db.Store, "par", "jafar", false)
		hotel          = fixtures.AddHotel(db.Store, "hilton", "tehran", 5, nil)
		room           = fixtures.AddRoom(db.Store, hotel.ID, 299.99, "large", true)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, room.ID, user.ID, 6, from, till)
		app            = fiber.New()
		bookingHandler = NewBookingHandler(*db.Store)
		route          = app.Group("/", middleware.JWTAuthentication(db.User))
	)
	route.Get("/:id", bookingHandler.HandleGetBooking)
	//test the page shows the actual booking
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)

	if resp.StatusCode != http.StatusOK {
		t.Fatal("expected status code 200 but got: ", resp.StatusCode)
	}

	var bookings *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	fmt.Println(bookings)
	if bookings.ID != booking.ID {
		t.Fatal("expected the same ID")
	}
	if bookings.NumPersons != booking.NumPersons {
		t.Fatal("expected the same number of persons")
	}
	if bookings.UserID != booking.UserID {
		t.Fatal("expected the same user but got", bookings.UserID)
	}
	if bookings.RoomID != booking.RoomID {
		t.Fatal("expected the same room but got a different one")
	}

	//test the unathenticated user
	route.Get("/:id", bookingHandler.HandleGetBooking)
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(noneAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatal("expected none 200 status for a unathenticated user")
	}
}
func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardrop(t)

	var (
		adminUser      = fixtures.AddUser(db.Store, "moji", "fery", true)
		user           = fixtures.AddUser(db.Store, "par", "jafar", false)
		hotel          = fixtures.AddHotel(db.Store, "hilton", "tehran", 5, nil)
		room           = fixtures.AddRoom(db.Store, hotel.ID, 299.99, "large", true)
		booking        = fixtures.AddBooking(db.Store, room.ID, adminUser.ID, 6, time.Now(), time.Now().AddDate(0, 0, 10))
		_              = booking
		app            = fiber.New()
		bookingHandler = NewBookingHandler(*db.Store)
		admin          = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
	)
	admin.Get("/", bookingHandler.HandelGetBookings)
	//test the page shows the correct bookings
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("expected 200 resp but got", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("Expected 1 booking but got %d", len(bookings))
	}
	if reflect.DeepEqual(booking, bookings[0]) {
		t.Fatal("expected bookings to be equal ")
	}

	//test none admin users could not access the page
	admin.Get("/", bookingHandler.HandelGetBookings)
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a none 200 status code but got %d", resp.StatusCode)
	}
}
