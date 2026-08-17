package utilities

import (
	"errors"
	"fmt"
	"log"
	"strings"

	mapStructure "github.com/go-viper/mapstructure/v2"
	fastHttp "github.com/valyala/fasthttp"
)

type requestClass struct {
	//radius float64 // Private field
}

func Request() *requestClass {
	return &requestClass{}
}

type AuthTokenInput struct {
	ApiID        string `bson:"api_id" json:"api_id"`
	UserID       string `bson:"user_id" json:"user_id"`
	LanguageCode string `bson:"lang_code" json:"lang_code"`
	Key          string `bson:"key" json:"key"`
}

type CustomerTokenInput struct {
	ApiID        string `bson:"api_id" json:"api_id"`
	CustomerID   string `bson:"customer_id" json:"customer_id"`
	LanguageCode string `bson:"lang_code" json:"lang_code"`
	Key          string `bson:"key" json:"key"`
}

func (r *requestClass) DecodeCustomerHeader(userAuth any) (CustomerTokenInput, error) {
	authDetail := make(map[string]string) //libRequest.CustomerTokenInput
	output := CustomerTokenInput{}
	err := mapStructure.Decode(userAuth, &authDetail)
	if err != nil {
		log.Println("Error decoding map to struct:", err)
		return output, err
	}
	output = CustomerTokenInput{
		ApiID:        authDetail["api_id"],
		CustomerID:   authDetail["customer_id"],
		LanguageCode: authDetail["lang_code"],
		Key:          authDetail["key"],
	}
	return output, nil
}

func (r *requestClass) DecodeUserHeader(userAuth any) (AuthTokenInput, error) {
	authDetail := make(map[string]string) //libRequest.AuthTokenInput
	output := AuthTokenInput{}
	err := mapStructure.Decode(userAuth, &authDetail)
	if err != nil {
		log.Println("Error decoding map to struct:", err)
		return output, err
	}
	output = AuthTokenInput{
		ApiID:        authDetail["api_id"],
		UserID:       authDetail["user_id"],
		LanguageCode: authDetail["lang_code"],
		Key:          authDetail["key"],
	}
	return output, nil
}

// GetBody extracts the body info from the request context
func (r *requestClass) PostBody(ctx *fastHttp.RequestCtx) []byte {
	bodyDetail, ok := ctx.UserValue("bodyDetail").(string)
	if !ok {
		return []byte("")
	}
	return []byte(bodyDetail)
}

func (r *requestClass) GetHeader(ctx *fastHttp.RequestCtx) (AuthTokenInput, error) {
	headerDetail, ok := ctx.UserValue("headerDetail").(AuthTokenInput)
	if !ok {
		return headerDetail, errors.New("unauthorized: header not found in context")
	}
	return headerDetail, nil
}

func (r *requestClass) GetBaseURL(ctx *fastHttp.RequestCtx, inUrl string) string {
	scheme := "http"
	if ctx.IsTLS() {
		scheme = "https"
	}
	output := fmt.Sprintf("%s://%s", scheme, ctx.Host())
	output = strings.Join([]string{output, inUrl}, "/")
	return output
}

func (r *requestClass) GetLangCode(ctx *fastHttp.RequestCtx) string {
	langCode, ok := ctx.UserValue("langCode").(string)
	if !ok {
		return "en"
	}
	return langCode
}

func (r *requestClass) GetServerUrl(ctx *fastHttp.RequestCtx) string {
	scheme := "http"
	if ctx.IsTLS() { ///file/:id
		scheme = "https"
	}
	serverURL := fmt.Sprintf("%s://%s", scheme, ctx.Host())
	return serverURL
}
func (r *requestClass) GetDownloadUrl(ctx *fastHttp.RequestCtx, inId string) string {
	serverURL := r.GetServerUrl(ctx) + "/file/" + inId
	return serverURL
}
