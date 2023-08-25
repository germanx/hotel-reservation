package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/germanx/hotel-reservation/api/middleware"
	"github.com/germanx/hotel-reservation/db/fixtures"
	"github.com/germanx/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		adminUser = fixtures.AddUser(db.Store, "admin", "admin", true)
		user      = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel     = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room      = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)

		d    = (60 * time.Second)
		from = time.Time(time.Now()).Round(d).UTC()
		till = from.AddDate(0, 0, 5).UTC()

		booking = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)

		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)
	// _ = booking

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response got: %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got: %d", len(bookings))
	}
	if !reflect.DeepEqual(booking, bookings[0]) {
		fmt.Printf("%+v\n", booking)
		fmt.Printf("%+v\n", bookings[0])
		t.Fatal("expected bookings to be equal")
	}

	// test non-admin cannot access the booking
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 response got: %d", resp.StatusCode)
	}
}

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		nonAuthUser = fixtures.AddUser(db.Store, "Jimmy", "Water", false)
		user        = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel       = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room        = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)

		d    = (60 * time.Second)
		from = time.Time(time.Now()).Round(d).UTC()
		till = from.AddDate(0, 0, 5).UTC()

		booking = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)

		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 code got %d", resp.StatusCode)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	// fmt.Println(bookingResp)

	if bookingResp.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, bookingResp.ID)
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookingResp.UserID)
	}

	// non auth user
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("non 200 code got %d", resp.StatusCode)
	}
}
