package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/germanx/hotel-reservation/db"
	"github.com/germanx/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore db.UserStore
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json: "user"`
	Token string      `json: "token"`
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredential(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResp{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

func NewAuthHandler(u db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: u,
	}
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	// fmt.Println("-- params:", params)

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// return fmt.Errorf("invalid credentials")
			return invalidCredential(c)
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(params.Password))
	if err != nil {
		return invalidCredential(c)
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return invalidCredential(c)
	}
	fmt.Println("authenticated ->", user)

	resp := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}
	return c.JSON(resp)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	// fmt.Println("-- secret:", secret)

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret")
	}
	return tokenStr
}
