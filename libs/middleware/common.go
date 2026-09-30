package middleware

import (
	"encoding/json"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	fastHttp "github.com/valyala/fasthttp"
)

type authRequest struct {
	CustomerId string `json:"customer_id" bson:"customer_id"`
	UserId     string `json:"user_id" bson:"user_id"`
}
type bodyRequest struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}
type CorsClass struct {
	origin  string
	methods string
}

var secretJwtKey string

func getSecretKey() []byte {
	temStr := strings.Split(time.Now().UTC().String(), " ")
	return []byte(secretJwtKey + temStr[0])
}
func verifyToken(tokenString string) (any, error) {
	secretKey := getSecretKey()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["data"], nil
	}
	return nil, nil
}

func extractBearerToken(ctx *fastHttp.RequestCtx) string {
	authHeader := ctx.Request.Header.Peek("Authorization")
	if len(authHeader) == 0 {
		return ""
	}
	authStr := string(authHeader)
	const prefix = "Bearer "
	if len(authStr) > len(prefix) && authStr[:len(prefix)] == prefix {
		return authStr[len(prefix):]
	}
	return ""
}
func extractHeader(ctx *fastHttp.RequestCtx, inKey string) string {
	keyHeader := ctx.Request.Header.Peek(inKey)
	if len(keyHeader) == 0 {
		return ""
	}
	return string(keyHeader)
}

func processAuthReq(ctx *fastHttp.RequestCtx, bodyRequest bodyRequest) (authRequest, bool) {
	authReq := authRequest{}
	authValue, err3 := verifyToken(bodyRequest.Value)
	statusOk := false
	if err3 == nil {
		temAuthValue, status := authValue.(string)
		if status == true {
			err4 := json.Unmarshal([]byte(temAuthValue), &authReq)
			if err4 == nil {
				ctx.Response.Header.Set("X-Customer-Id", authReq.CustomerId)
				ctx.Response.Header.Set("X-User-Id", authReq.UserId)
				statusOk = true
			}
		}
	}
	return authReq, statusOk
}

func processBodyReq(ctx *fastHttp.RequestCtx, bodyRequest bodyRequest) (string, bool) {
	bodyValue, err2 := verifyToken(bodyRequest.Key)
	if err2 == nil {
		temBodyValue, status := bodyValue.(string)
		if status == true {
			ctx.Request.SetBodyString(temBodyValue)
		}
		return temBodyValue, status
	}
	return "", false
}
