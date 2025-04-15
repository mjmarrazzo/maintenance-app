package validation

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mjmarrazzo/maintenance-app/internal/responses"
)

func BindPathParams(c echo.Context, i interface{}) error {
	if err := bindPathParamsOnly(c, i); err != nil {
		return handlePathParamError(i, err)
	}

	if err := ValidateStruct(i); err != nil {
		return handlePathParamError(i, err)
	}

	return nil
}

func bindPathParamsOnly(c echo.Context, i interface{}) error {
	typ := reflect.TypeOf(i).Elem()
	val := reflect.ValueOf(i).Elem()

	for i := range typ.NumField() {
		field := typ.Field(i)
		paramTag := field.Tag.Get("param")

		if paramTag == "" {
			continue
		}

		paramValue := c.Param(paramTag)
		if paramValue == "" {
			continue
		}

		fieldValue := val.Field(i)

		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(paramValue, 10, 64)
			if err != nil {
				return &PathParamError{
					ParamName:    paramTag,
					ParamValue:   paramValue,
					ExpectedType: "integer",
				}
			}
			fieldValue.SetInt(intVal)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uintVal, err := strconv.ParseUint(paramValue, 10, 64)
			if err != nil {
				return &PathParamError{
					ParamName:    paramTag,
					ParamValue:   paramValue,
					ExpectedType: "unsigned integer",
				}
			}
			fieldValue.SetUint(uintVal)

		case reflect.Float32, reflect.Float64:
			floatVal, err := strconv.ParseFloat(paramValue, 64)
			if err != nil {
				return &PathParamError{
					ParamName:    paramTag,
					ParamValue:   paramValue,
					ExpectedType: "float",
				}
			}
			fieldValue.SetFloat(floatVal)

		case reflect.Bool:
			boolVal, err := strconv.ParseBool(paramValue)
			if err != nil {
				return &PathParamError{
					ParamName:    paramTag,
					ParamValue:   paramValue,
					ExpectedType: "boolean",
				}
			}
			fieldValue.SetBool(boolVal)

		case reflect.String:
			fieldValue.SetString(paramValue)
		}
	}

	return nil
}

type PathParamError struct {
	ParamName    string
	ParamValue   interface{}
	ExpectedType string
}

func (e *PathParamError) Error() string {
	return fmt.Sprintf("failed to parse path parameter '%s' with value '%s' as %s",
		e.ParamName, e.ParamValue, e.ExpectedType)
}

func handlePathParamError(i interface{}, err error) error {
	validationResult := &responses.ValidationErrors{
		Parameters: []string{},
		Violations: []*responses.ViolationsDetail{},
	}

	switch typedErr := err.(type) {
	case *PathParamError:
		validationResult.Parameters = append(validationResult.Parameters, typedErr.ParamName)
		validationResult.Violations = append(validationResult.Violations, &responses.ViolationsDetail{
			Name:    typedErr.ParamName,
			Message: GetExpectedTypeErrorMessage(typedErr.ExpectedType, nil),
		})

	case validator.ValidationErrors:
		for _, e := range typedErr {
			fieldName := e.Field()
			paramName := getParamNameFromField(i, fieldName)

			validationResult.Parameters = append(validationResult.Parameters, paramName)
			validationResult.Violations = append(validationResult.Violations, &responses.ViolationsDetail{
				Name:    paramName,
				Message: GetErrorMessage(e),
			})
		}

	default:
		validationResult.Violations = append(validationResult.Violations, &responses.ViolationsDetail{
			Name:    "unknown",
			Message: fmt.Sprintf("Path parameter error: %v", err),
		})
	}

	return NewValidationError("Validation failed", validationResult)
}

func getParamNameFromField(i interface{}, fieldName string) string {
	typ := reflect.TypeOf(i).Elem()
	field, found := typ.FieldByName(fieldName)
	if !found {
		return strings.ToLower(fieldName)
	}

	paramTag := field.Tag.Get("param")
	if paramTag != "" {
		return paramTag
	}

	return strings.ToLower(fieldName)
}
