package utils

import (
	"time"

	"github.com/armiariyan/synapsis/internal/config"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	ID    string
	Email string
	jwt.StandardClaims
}

func JwtSign(id string, email string) (token string, exp int64, err error) {
	secret := []byte(config.GetString("jwt"))
	exp = time.Now().Add(7 * 24 * time.Hour).Unix()
	claims := &Claims{
		id,
		email,
		jwt.StandardClaims{
			Issuer:    "DOMPET_KILAT",
			ExpiresAt: exp, // * 7 days
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	return
}

func JwtVerify(token string) (claims *Claims, err error) {
	secret := []byte(config.GetString("jwt"))
	claims = &Claims{}
	_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	return
}
