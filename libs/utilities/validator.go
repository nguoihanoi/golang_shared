package utilities

import (
	"encoding/json"

	mainValidator "github.com/go-playground/validator/v10"
	fastHttp "github.com/valyala/fasthttp"
)

var statusCreate bool = false
var validate *mainValidator.Validate

func GetValidate() *mainValidator.Validate {
	if !statusCreate {
		validate = mainValidator.New(mainValidator.WithRequiredStructEnabled())
		statusCreate = true
	}
	return validate
}

func Validate(ctx *fastHttp.RequestCtx, regRequest any) (err error) {
	//Todo: get struct input
	if err = json.Unmarshal(ctx.PostBody(), regRequest); err != nil {
		return err
	}
	//Todo: validate struct input
	mainValidate := GetValidate()
	err = mainValidate.Struct(regRequest)
	return err
}

func Validate2(ctx *fastHttp.RequestCtx) (regRequest any, err error) {
	//Todo: get struct input
	if err = json.Unmarshal(ctx.PostBody(), regRequest); err != nil {
		return nil, err
	}
	//Todo: validate struct input
	mainValidate := GetValidate()
	err = mainValidate.Struct(regRequest)
	return regRequest, err
}

func ValidateLangValue(inMap map[string]string, isEmpty bool, langCodes []string) (results []string, status int) {
	status = 0
	keys := GetMapKeys(inMap)
	arrResult := SymmetricDifference(keys, langCodes)
	if len(arrResult) > 0 {
		status = 1
	}
	for i := range keys {
		temValue := String().Trim(inMap[keys[i]])
		if temValue == "" {
			results = append(results, keys[i])
			break
		}
	}
	if len(results) > 0 {
		status = 2
	}
	return results, status
}
