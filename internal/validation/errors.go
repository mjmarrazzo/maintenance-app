package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mjmarrazzo/maintenance-app/internal/responses"
)

func NewValidationError(message string, validationErrors *responses.ValidationErrors) error {
	if validationErrors == nil {
		return responses.NewValidationError(message, nil, nil)
	}

	return responses.NewValidationError(message, validationErrors.Parameters, validationErrors.Violations)
}

func ValidateStruct(i interface{}) error {
	validate := GetValidator()
	return validate.Validator.Struct(i)
}

func GetValidationError(err error) (*responses.ValidationError, bool) {
	valErr, ok := err.(responses.ValidationError)
	return &valErr, ok
}

func HandleValidationErrors(err error) error {
	if err == nil {
		return nil
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	validationResult := &responses.ValidationErrors{
		Parameters: make([]string, 0),
		Violations: make([]*responses.ViolationsDetail, 0),
	}

	for _, err := range errs {
		field := err.Field()
		namespace := err.Namespace()

		if idx := findArrayIndex(namespace); idx != "" {
			field = field + idx
		} else if parent := findParentField(namespace); parent != "" {
			field = parent + "." + field
		}

		validationResult.Parameters = append(validationResult.Parameters, field)

		validationResult.Violations = append(validationResult.Violations, &responses.ViolationsDetail{
			Name:    field,
			Message: GetErrorMessage(err),
		})
	}

	return NewValidationError("Validation failed", validationResult)
}

func findArrayIndex(namespace string) string {
	re := regexp.MustCompile(`\[(\d+)\]`)
	matches := re.FindStringSubmatch(namespace)

	if len(matches) > 1 {
		return "[" + matches[1] + "]"
	}

	return ""
}

func findParentField(namespace string) string {
	parts := strings.Split(namespace, ".")

	if len(parts) < 3 {
		return ""
	}

	return parts[1]
}

func GetErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Should be at least %s", err.Param())
	case "max":
		return fmt.Sprintf("Should be at most %s", err.Param())
	case "ascii":
		return "Should contain only ASCII characters"
	case "numericstring":
		return "Should be a numeric string"
	case "gt":
		if err.Param() == "0" {
			return "Should be a positive number"
		}
		return fmt.Sprintf("Should be greater than %s", err.Param())
	case "gte":
		if err.Param() == "1" {
			return "Should be a positive number or zero"
		}
		return fmt.Sprintf("Should be greater than or equal to %s", err.Param())
	case "lt":
		if err.Param() == "0" {
			return "Should be a negative number"
		}
		return fmt.Sprintf("Should be less than %s", err.Param())
	case "lte":
		if err.Param() == "-1" {
			return "Should be a negative number or zero"
		}
		return fmt.Sprintf("Should be less than or equal to %s", err.Param())
	default:
		return fmt.Sprintf("Failed validation: %s", err.Tag())
	}
}

func GetExpectedTypeErrorMessage(expectedType string, value interface{}) string {
	if value == nil {
		return fmt.Sprintf("Expected type <%s>", expectedType)
	}
	return fmt.Sprintf("Expected type <%s>, but got <%s>", expectedType, value)
}
