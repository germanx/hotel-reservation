package api

import (
	"net/http"

	"github.com/germanx/hotel-reservation/db"
	"github.com/germanx/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
	list, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(list)
}

// TODO: this needs to be user authorized!
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	item, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return err
	}
	if item.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}
	return c.JSON(item)
}
