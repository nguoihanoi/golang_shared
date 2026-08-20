package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	fastHttp "github.com/valyala/fasthttp"
)

// var languageCodeCache *libCache.Cache
type ResponseClass struct {
	//radius float64 // Private field
	cache *libCache.Cache
	key   string
}

func Response(inRedisClient *libCache.Cache, inKey string) *ResponseClass {
	return &ResponseClass{
		cache: inRedisClient,
		key:   inKey + ":",
	}
}
func (r *ResponseClass) getCodeByKey(inKey, inLanguageCode string) string {
	cacheKey := r.key + inLanguageCode
	if val, ok := r.cache.HGet(cacheKey, inKey).(string); ok && val != "" {
		return val
	}
	r.cache.HSet(cacheKey, inKey, inKey)
	return inKey
}

// Message detail
func (r *ResponseClass) Message(status bool, message string, inStatusCode int) map[string]any {
	return map[string]any{"status": status, "message": message, "statusCode": inStatusCode}
}

// Send detail
func (r *ResponseClass) Send(ctx *fastHttp.RequestCtx, data map[string]any) {
	ctx.Request.Header.Set("Accept-Encoding", "gzip, deflate, sdhc")
	ctx.Request.Header.Add("Content-Encoding", "gzip")
	ctx.Response.Header.SetCanonical([]byte("Content-Type"), []byte("application/json"))
	err := json.NewEncoder(ctx).Encode(data)
	if err != nil {
		log.Println(err)
	}
}

type ContentResponseOutput struct {
	Status     bool   `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}
type ResponseOutput struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (r *ResponseClass) GetOutput(inStatus bool, inMessage string, inStatusCode int) (output ContentResponseOutput) {
	output.Status = inStatus
	output.Message = inMessage
	output.StatusCode = inStatusCode
	output.Data = nil
	return output
}
func (r *ResponseClass) SendOutput(ctx *fastHttp.RequestCtx, inResponse ContentResponseOutput) {
	inResponse.Message = r.getCodeByKey(inResponse.Message, "vi")
	responseOutput := ResponseOutput{
		Status:  inResponse.Status,
		Message: inResponse.Message,
		Data:    inResponse.Data,
	}
	if inResponse.Status {
		inResponse.StatusCode = 200
	}
	jsonResponse, _ := json.Marshal(responseOutput)
	ctx.Request.Header.Set("Accept-Encoding", "gzip, deflate, sdhc")
	ctx.Request.Header.Add("Content-Encoding", "gzip")
	ctx.SetStatusCode(inResponse.StatusCode)
	ctx.SetContentType("application/json")
	ctx.SetBody(jsonResponse)
}

func (r *ResponseClass) SendError(ctx *fastHttp.RequestCtx, inMessage string, inError libProcess.E, inStatusCode int) {
	log.Println(inError)
	inMessage = r.getCodeByKey(inMessage, "vi")
	resp := r.GetOutput(false, inMessage, inStatusCode)
	resp.Data = fmt.Sprint(inError)
	r.SendOutput(ctx, resp)
}

// GetUserIDFromContext extracts the user ID from the request context
// This function relies on the AuthMiddleware having previously validated the token
// and added the user ID to the context.
func (r *ResponseClass) GetUserIDFromContext(ctx *fastHttp.RequestCtx) (string, error) {
	userID, ok := ctx.UserValue("userId").(string)
	if !ok || userID == "" {
		return "", errors.New("unauthorized: user ID not found in context")
	}
	return userID, nil
}
