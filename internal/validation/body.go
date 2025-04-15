package validation

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mjmarrazzo/maintenance-app/internal/responses"
)

func BindBody(e echo.Context, i interface{}) error {
	// if err := validateNoExtraFields(e, i); err != nil {
	// 	e.Logger().Print("failed to validate extra fields", err)
	// 	return err
	// }

	if err := e.Bind(i); err != nil {
		return HandleBindError(e, err, reflect.TypeOf(i))
	}

	if err := ValidateStruct(i); err != nil {
		return HandleValidationErrors(err)
	}

	return nil
}

func HandleBindError(e echo.Context, err error, targetType reflect.Type) error {
	validationResult := &responses.ValidationErrors{
		Parameters: []string{},
		Violations: []*responses.ViolationsDetail{},
	}

	switch typedErr := err.(type) {
	case *json.UnmarshalTypeError:
		fieldName := typedErr.Field
		value := typedErr.Value
		expectedType := typedErr.Type.String()

		validationResult.Parameters = append(validationResult.Parameters, fieldName)
		validationResult.Violations = append(validationResult.Violations, &responses.ViolationsDetail{
			Name:    fieldName,
			Message: GetExpectedTypeErrorMessage(expectedType, value),
		})
	case *json.SyntaxError:
		return NewValidationError("Invalid JSON syntax", validationResult)
	case *echo.HTTPError:
		ie := typedErr.Internal
		if ie != nil {
			log.Printf("Recursively handling internal type: %T\n", ie)
			return HandleBindError(e, ie, targetType)
		}
	case *responses.ValidationError:
		validationResult.Parameters = append(validationResult.Parameters, typedErr.Parameters...)
		validationResult.Violations = append(validationResult.Violations, typedErr.Violations...)

	default:
		log.Printf("Unknown error type: %T, message: %v\n", err, err)
	}

	return NewValidationError("Validation failed", validationResult)
}

func validateNoExtraFields(e echo.Context, obj interface{}) error {
	body, err := io.ReadAll(e.Request().Body)
	if err != nil {
		return err
	}

	e.Request().Body = io.NopCloser(bytes.NewReader(body))

	var requestMap map[string]interface{}
	if err := json.Unmarshal(body, &requestMap); err != nil {
		return NewValidationError("Invalid JSON syntax", nil)
	}

	structFields := make(map[string]bool)
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag != "" {
			if commaIdx := strings.Index(tag, ","); commaIdx != -1 {
				tag = tag[:commaIdx]
			}
		}
		structFields[tag] = true
	}

	var extraFields []string
	for key := range requestMap {
		if !structFields[key] {
			extraFields = append(extraFields, key)
		}
	}

	if len(extraFields) > 0 {
		validationResult := &responses.ValidationErrors{
			Parameters: make([]string, 0),
			Violations: make([]*responses.ViolationsDetail, 0),
		}
		for _, field := range extraFields {
			validationResult.Parameters = append(validationResult.Parameters, field)
			validationResult.Violations = append(validationResult.Violations, &responses.ViolationsDetail{
				Name:    field,
				Message: "Unexpected field",
			})
		}

		return NewValidationError("Validation failed", validationResult)
	}

	return nil
}
