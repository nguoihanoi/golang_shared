package crypto

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("ThanhNv")
var prefixToken string

type JwtClass struct {
	//radius float64 // Private field
}

func (c *CryptoClass) JWT() *JwtClass {
	return &JwtClass{}
}

func (j *JwtClass) SetToken(inJwtToken string, inPrefixToken string) {
	secretKey = []byte(inJwtToken)
	prefixToken = inPrefixToken
}
func (j *JwtClass) CreateToken(inData any) (string, time.Time, error) {
	nextTime := time.Now().Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"data": inData,
			"exp":  nextTime.Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", nextTime, err
	}

	return tokenString, nextTime, nil
}

func (j *JwtClass) VerifyToken(tokenString string) (any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["data"], nil
	}
	return nil, nil
}
