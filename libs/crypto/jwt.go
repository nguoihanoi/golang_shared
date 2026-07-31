package crypto

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JwtClass struct {
	secretKey string
}

func JWT(inJwtToken string) *JwtClass {
	return &JwtClass{secretKey: inJwtToken}
}

func (j *JwtClass) CreateToken(inData any) (string, time.Time, error) {
	nextTime := time.Now().Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"data": inData,
			"exp":  nextTime.Unix(),
		})

	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", nextTime, err
	}

	return tokenString, nextTime, nil
}

func (j *JwtClass) VerifyToken(tokenString string) (any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["data"], nil
	}
	return nil, nil
}
