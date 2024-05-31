package jwt

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/cast"
)

type contextKey struct{}

type Payload struct {
	UserID int    `json:"userId"`
	Email  string `json:"email"`
	Exp    int64  `json:"exp"`
	Role   string `json:"role"`
}

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenNotFound      = errors.New("token not found")
	ErrInvalidClaimFormat = errors.New("invalid claims format")
)

func SavePayloadToContext(parent context.Context, payload Payload) context.Context {
	return context.WithValue(parent, contextKey{}, payload)
}

func GetPayloadFromContext(ctx context.Context) Payload {
	payload, ok := ctx.Value(contextKey{}).(Payload)
	if ok {
		return payload
	}

	return Payload{}
}

func Validate(tokenString string, secretKey []byte) (Payload, error) {
	if tokenString == "" {
		return Payload{}, ErrTokenNotFound
	}

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key for validation
		return secretKey, nil
	})

	// Got error when parsing JWT token
	if err != nil {
		return Payload{}, err
	}

	// Token is invalid
	if !token.Valid {
		return Payload{}, ErrInvalidToken
	}

	// Access the claims as a map
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Payload{}, ErrInvalidClaimFormat
	}

	payload := Payload{
		UserID: cast.ToInt(claims["userId"]),
		Email:  cast.ToString(claims["email"]),
		Exp:    cast.ToInt64(claims["exp"]),
	}

	// Return the token payload
	return payload, nil
}
