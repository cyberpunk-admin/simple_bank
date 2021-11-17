package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/simplebank/db/util"
)

var validateCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupport(currency)
	}
	return false
}
