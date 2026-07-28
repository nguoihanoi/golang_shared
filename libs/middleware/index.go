package middleware

import (
	"log"
	"time"

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
