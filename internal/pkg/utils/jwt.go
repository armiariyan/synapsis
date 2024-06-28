package utils

import (
	"encoding/json"
	"time"

	"github.com/armiariyan/synapsis/internal/config"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Data any `json:"data"`
	jwt.StandardClaims
}

type JWTClaimsData struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

func JwtSign(data any) (signedString string, exp int64, err error) {
	secret := []byte(config.GetString("jwt.secret"))
	bt, err := json.Marshal(data)
	if err != nil {
		return
	}
	exp = time.Now().Add(7 * 24 * time.Hour).Unix() // * 7 days
	encryptedData, err := EncryptAES256CBC(string(bt), config.GetString("aes256cbc.encrypt.key"), config.GetString("aes256cbc.encrypt.iv"))
	if err != nil {
		return
	}
	claims := &Claims{
		encryptedData,
		jwt.StandardClaims{
			Issuer:    config.GetString("jwt.issuer"),
			ExpiresAt: exp,
		},
	}
	signedString, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	return
}

func JwtVerify(token string) (claims *Claims, err error) {
	secret := []byte(config.GetString("jwt.secret"))
	claims = &Claims{}
	_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return
	}
	decryptedString, err := DecryptAES256CBC(claims.Data.(string), config.GetString("aes256cbc.encrypt.key"), config.GetString("aes256cbc.encrypt.iv")) // * should be string
	if err != nil {
		return
	}
	var data JWTClaimsData
	err = json.Unmarshal([]byte(decryptedString), &data)
	if err != nil {
		return
	}

	claims.Data = data

	return
}
