package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type masterClaims struct {
	masterId    string
	masterEmail string
	jwt.RegisteredClaims
}

func GenerateToken(masterId, masterEmail string) (string, error) {
	claims := masterClaims{
		masterId,
		masterEmail,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte("thisisasecretstring"))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
