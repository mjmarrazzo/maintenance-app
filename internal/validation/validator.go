package validation

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidatorInstance struct {
	Validator  *validator.Validate
	Translator ut.Translator
}

var once sync.Once
var instance *ValidatorInstance

func GetValidator() *ValidatorInstance {
	once.Do(func() {
		v := validator.New(validator.WithRequiredStructEnabled())

		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return fld.Name
			}
			return name
		})
		RegisterNumericStringValidator(v)

		instance = &ValidatorInstance{
			Validator:  v,
			Translator: nil,
		}
	})

	return instance
}

func RegisterNumericStringValidator(v *validator.Validate) {
	numericRegex := regexp.MustCompile(`^-?[0-9]+(\.[0-9]+)?$`)
	err := v.RegisterValidation("numericstring", func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() != reflect.String {
			return true
		}

		fieldValue := fl.Field().String()

		if fieldValue == "" {
			return true
		}

		if !numericRegex.MatchString(fieldValue) {
			return false
		}

		tag := fl.GetTag()

		additionalTags := tag
		if strings.Contains(additionalTags, "numericstring") {
			additionalTags = strings.ReplaceAll(additionalTags, "numericstring,", "")
			additionalTags = strings.ReplaceAll(additionalTags, ",numericstring", "")
			additionalTags = strings.ReplaceAll(additionalTags, "numericstring", "")
		}

		if additionalTags == "" {
			return true
		}

		if val, err := strconv.ParseInt(fieldValue, 10, 64); err == nil {
			if err := v.Var(val, additionalTags); err == nil {
				return true
			}
		}

		if val, err := strconv.ParseFloat(fieldValue, 64); err == nil {
			return v.Var(val, additionalTags) == nil
		}

		return false
	})

	if err != nil {
		panic(err)
	}
}
