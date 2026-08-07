package middleware

import (
	"log"
	"time"

	libCrypto "github.com/nguoihanoi/golang_shared/libs/crypto"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	fastHttp "github.com/valyala/fasthttp"
)

type Middleware func(fastHttp.RequestHandler) fastHttp.RequestHandler

// MwCompress dteail...
func MwCompress(h fastHttp.RequestHandler) fastHttp.RequestHandler {
	return fastHttp.CompressHandler(func(ctx *fastHttp.RequestCtx) {
		h(ctx)
	})
}

// Logging logs all requests with its path and the time it took to process
func Logging() Middleware {
	// Create a new Middleware
	return func(h fastHttp.RequestHandler) fastHttp.RequestHandler {
		return func(ctx *fastHttp.RequestCtx) {
			// Do middleware things
			start := time.Now()
			defer func() {
				log.Println(string(ctx.Path()), time.Since(start))
			}()
			h(ctx)
		}
	}
}
func Compress() Middleware {
	// Create a new Middleware
	return func(h fastHttp.RequestHandler) fastHttp.RequestHandler {
		return fastHttp.CompressHandler(func(ctx *fastHttp.RequestCtx) {
			h(ctx)
		})
	}
}

// Post applies middlewares to fastHttp.RequestHandler
func Post(h fastHttp.RequestHandler) fastHttp.RequestHandler {
	middlewares := []Middleware{}
	middlewares = append(middlewares, Logging())
	middlewares = append(middlewares, Compress())
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

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

var libJwt *libCrypto.JwtClass

func Init(inOrigin string, inMethod string, inToken string) *CorsClass {
	libJwt = libCrypto.JWT(inToken)
	return &CorsClass{
		origin:  inOrigin,
		methods: inMethod,
	}
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

func (c *CorsClass) CorsMiddleware(next fastHttp.RequestHandler) fastHttp.RequestHandler {
	output := func(ctx *fastHttp.RequestCtx) {
		// Set CORS headers
		start := time.Now()
		ctx.Response.Header.Set("Access-Control-Allow-Origin", c.origin)
		ctx.Response.Header.Set("Access-Control-Expose-Headers", "Authorization")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", c.methods)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, Cache-Control")
		// Handle preflight (OPTIONS) requests
		if string(ctx.Method()) == "OPTIONS" {
			ctx.SetStatusCode(fastHttp.StatusOK)
			return
		}
		//var token = extractBearerToken(ctx)
		var apiKey = extractHeader(ctx, "X-Api-Key")
		if apiKey == "1" {
			var bodyRequest bodyRequest
			err := libUtilities.Validate(ctx, &bodyRequest)
			if err == nil {
				//Todo: verify auth
				authValue, err3 := libJwt.VerifyToken(bodyRequest.Key)
				if err3 == nil {
					temBodyValue, status := authValue.(authRequest)
					if status == true {
						log.Println(temBodyValue)
						ctx.Response.Header.Set("X-Customer-Id", temBodyValue.CustomerId)
						ctx.Response.Header.Set("X-User-Id", temBodyValue.UserId)
					} else {
						ctx.SetStatusCode(fastHttp.StatusForbidden)
						return
					}
				} else {
					ctx.SetStatusCode(fastHttp.StatusForbidden)
					return
				}
				//Todo: verify body
				bodyValue, err2 := libJwt.VerifyToken(bodyRequest.Key)
				if err2 == nil {
					temBodyValue, status := bodyValue.(string)
					if status == true {
						ctx.Request.SetBodyString(temBodyValue)
					} else {
						ctx.SetStatusCode(fastHttp.StatusForbidden)
						return
					}
				} else {
					ctx.SetStatusCode(fastHttp.StatusForbidden)
					return
				}
			} else {
				ctx.SetStatusCode(fastHttp.StatusForbidden)
				return
			}
		}
		// Do middleware things
		defer func() {
			log.Println(string(ctx.Path()), time.Since(start))
		}()
		// Call the next handler
		next(ctx)
	}
	return fastHttp.CompressHandler(output)
}
