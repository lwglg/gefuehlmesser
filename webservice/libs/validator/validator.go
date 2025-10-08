package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"

	s "webservice/libs/sentiment"
)

type ErrResponse struct {
	Errors []string `json:"errors"`
}

func New() *validator.Validate {
	validate := validator.New()

	// Using the names which have been specified for JSON representations of structs, rather than normal Go field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	validate.RegisterValidation("id_field", isIdFieldValid)
	validate.RegisterValidation("datetime", isTimestampValid)

	return validate
}

func ToErrResponse(err error) *ErrResponse {
	fieldErrors, ok := err.(validator.ValidationErrors)

	if ok {
		resp := ErrResponse{
			Errors: make([]string, len(fieldErrors)),
		}

		for i, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				resp.Errors[i] = fmt.Sprintf("O campo '%s' é obrigatório", err.Field())
			case "max":
				resp.Errors[i] = fmt.Sprintf("O campo '%s' deve ter um tamanho de no máximo %s caracteres", err.Field(), err.Param())
			case "id_field":
				resp.Errors[i] = fmt.Sprintf("O campo '%s' deve seguir o padrão '[user|msg|perf]_XXX'", err.Field())
			case "datetime":
				resp.Errors[i] = fmt.Sprintf("O campo '%s' deve seguir o formato RFC3339 'YYYY-MM-DDThh:mm:ssZ'", err.Field())
			default:
				resp.Errors[i] = fmt.Sprintf("Ocorreu um erro ao realizar a validação de tag '%s' no campo '%s'", err.Tag(), err.Field())
			}
		}

		return &resp
	}

	return nil
}

func isRegexFieldValid(fl validator.FieldLevel, regexString string) bool {
	reg := regexp.MustCompile(regexString)

	return reg.MatchString(fl.Field().String())
}

func isIdFieldValid(fl validator.FieldLevel) bool {
	return isRegexFieldValid(fl, s.IdRegexString)
}

func isTimestampValid(fl validator.FieldLevel) bool {
	return isRegexFieldValid(fl, s.TimeStampRFC3339RegexString)
}
