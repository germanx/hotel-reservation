package api

import (
	"github.com/germanx/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(s *db.Store) *BookingHandler {
	return &BookingHandler{
		store: s,
	}
}

// TODO: this needs to be admin authorized!
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	list, err := h.store.Booking.GetBookings(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(list)
}

// TODO: this needs to be user authorized!
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	list, err := h.store.Booking.GetBookings(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(list)
}
