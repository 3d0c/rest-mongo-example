package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/teal-seagull/lyre-be-v4/pkg/config"
)

// AuthClaims embeddes StandardClaims and stores UserID
type AuthClaims struct {
	UserID string
	jwt.StandardClaims
}

// GetUserID is a userID getter. Do not use UserID field directly.
// It should be exported because of JSON marshalling.
func (a *AuthClaims) GetUserID() (int, error) {
	return strconv.Atoi(a.UserID)
}

// CreateToken generates token string
func CreateToken(userID string) (string, error) {
	var (
		secret []byte = []byte(config.TheConfig().Server.JWTSecret)
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
			Issuer:    "lyre-be-v4",
		},
	})

	tokenString, err := token.SignedString(secret)

	return tokenString, err
}

// VerifyToken parses and validates token string end returns userID from MapClaims
func VerifyToken(tokenString string) (*AuthClaims, error) {
	var (
		secret []byte = []byte(config.TheConfig().Server.JWTSecret)
		ok     bool
	)

	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token '%s'", tokenString)
	}

	if _, ok = token.Claims.(*AuthClaims); !ok {
		return nil, fmt.Errorf("error casting Claims to AuthClaims")
	}

	return token.Claims.(*AuthClaims), nil
}
