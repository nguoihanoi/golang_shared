package middleware

import (
	"log"
	"sync"
	"time"

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

func Init(inOrigin string, inMethod string, inToken string) *CorsClass {
	secretJwtKey = inToken
	return &CorsClass{
		origin:  inOrigin,
		methods: inMethod,
	}
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
				var (
					authReq      *authRequest
					temBodyValue string
				)
				statusOk := true
				var wg sync.WaitGroup
				wg.Add(2)
				// 1. Luồng 1: Xử lý Xác thực (Auth)
				go func() {
					defer wg.Done()
					res, ok := processAuthReq(ctx, bodyRequest)
					if !ok {
						statusOk = false
					}
					authReq = &res
				}()
				// 2. Luồng 2: Xử lý Verify Body
				go func() {
					defer wg.Done()
					bodyVal, ok := processBodyReq(ctx, bodyRequest)
					if !ok {
						statusOk = false
					}
					temBodyValue = bodyVal
				}()
				wg.Wait()
				if statusOk == false {
					// Nếu 1 trong 2 luồng thất bại -> Trả về 403 Forbidden
					ctx.SetStatusCode(fastHttp.StatusForbidden)
					return
				}

				// 4. Ghi thông tin vào ctx trên Goroutine CHÍNH (Thread-safe)
				ctx.Response.Header.Set("X-Customer-Id", authReq.CustomerId)
				ctx.Response.Header.Set("X-User-Id", authReq.UserId)
				ctx.Request.SetBodyString(temBodyValue)
			} else {
				ctx.SetStatusCode(fastHttp.StatusForbidden)
				return
			}
		}
		log.Println("done")
		// Do middleware things
		defer func() {
			log.Println(string(ctx.Path()), time.Since(start))
		}()
		// Call the next handler
		next(ctx)
	}
	return fastHttp.CompressHandler(output)
}
