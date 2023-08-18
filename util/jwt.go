package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"simple_douyin/config"
)

var jwtSecret = []byte(config.JWTSECRET)

type Claims struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
	jwt.StandardClaims
}

func GenerateToken(id int64, username string) (string, error) {
	expireTime := time.Now().Add(config.TokenExpireDuration)

	claims := Claims{
		username,
		id,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gs",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
