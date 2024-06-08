package api

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mojtabafarzaneh/booking_room/db"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, ok := ctx.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return ErrUnAuthorized()
		}
		_, err := validatedToken(token[0])

		if err != nil {
			return err
		}
		claims, err := validatedToken(token[0])
		if err != nil {
			return ErrUnAuthorized()
		}

		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(ctx.Context(), userID)
		if err != nil {
			return err
		}
		ctx.Context().SetUserValue("user", user)
		return ctx.Next()
	}
}

func validatedToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signature method", token.Header["alg"])
			return nil, ErrUnAuthorized()
		}

		secret := os.Getenv("JWT_SECRET")
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, ErrUnAuthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		fmt.Println(claims["foo"], claims["nbf"])
	}
	return claims, nil
}
