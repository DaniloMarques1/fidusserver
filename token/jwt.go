package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type masterClaims struct {
	MasterId    string `json:"master_id"`
	MasterEmail string `json:"master_email"`
	jwt.RegisteredClaims
}

func GenerateToken(masterId, masterEmail string) (string, int64, error) {
	expiresAt := time.Now().Add(5 * time.Minute)
	claims := &masterClaims{
		MasterId:    masterId,
		MasterEmail: masterEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	jwtSecretKey := os.Getenv("JWT_KEY")
	signedToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", 0, err
	}
	return signedToken, expiresAt.Unix(), nil
}

func ParseToken(tokenStr string) (string, error) {
	jwtSecretKey := os.Getenv("JWT_KEY")
	token, err := jwt.ParseWithClaims(tokenStr, &masterClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return "", err
	}
	claims := token.Claims.(*masterClaims)
	return claims.MasterId, nil
}
