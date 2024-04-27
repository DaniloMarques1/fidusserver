package token

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type masterClaims struct {
	MasterId    string
	MasterEmail string
	jwt.RegisteredClaims
}

func GenerateToken(masterId, masterEmail string) (string, error) {
	claims := &masterClaims{
		MasterId:    masterId,
		MasterEmail: masterEmail,
		RegisteredClaims: jwt.RegisteredClaims{
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

func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &masterClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("thisisasecretstring"), nil
	})
	if err != nil {
		return "", err
	}
	claims := token.Claims.(*masterClaims)
	log.Printf("Claims = %v\n", claims)
	return claims.MasterId, nil
}
