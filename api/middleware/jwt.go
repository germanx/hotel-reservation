package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/germanx/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// fmt.Println("-- JWT authing")
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			fmt.Println("token not present in the header")
			return fmt.Errorf("unauthorized")
		}
		// fmt.Println("-- token:", token)

		claims, err := validateToken(token)
		if err != nil {
			return err
		}
		// check token expiration
		expiresF := claims["expires"].(float64)
		expires := int64(expiresF)
		// expires := claims["expires"].(int64)
		if time.Now().Unix() > expires {
			return fmt.Errorf("token expires")
		}

		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}
		// set user to context
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		// fmt.Println("SECRET:", secret)
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println(claims)
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil
}
