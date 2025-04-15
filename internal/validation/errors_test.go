package validation

import (
	"fmt"
	"reflect"
	"testing"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/mjmarrazzo/maintenance-app/internal/responses"
)

type mockFieldError struct {
	tag   string
	param string
}

func (m mockFieldError) Tag() string                       { return m.tag }
func (m mockFieldError) ActualTag() string                 { return m.tag }
func (m mockFieldError) Namespace() string                 { return "" }
func (m mockFieldError) StructNamespace() string           { return "" }
func (m mockFieldError) Field() string                     { return "" }
func (m mockFieldError) StructField() string               { return "" }
func (m mockFieldError) Value() interface{}                { return nil }
func (m mockFieldError) Param() string                     { return m.param }
func (m mockFieldError) Kind() reflect.Kind                { return reflect.String }
func (m mockFieldError) Type() reflect.Type                { return reflect.TypeOf("") }
func (m mockFieldError) Error() string                     { return "" }
func (m mockFieldError) Translate(ut ut.Translator) string { return "" }

func TestNewValidationError(t *testing.T) {
	tests := []struct {
		name            string
		message         string
		validationErrs  *responses.ValidationErrors
		expectedCode    string
		expectedMessage string
	}{
		{
			name:            "nil validation errors",
			message:         "Test message",
			validationErrs:  nil,
			expectedCode:    "INVALID_FORMAT",
			expectedMessage: "Test message",
		},
		{
			name:    "with validation errors",
			message: "Test message",
			validationErrs: &responses.ValidationErrors{
				Parameters: []string{"field1"},
				Violations: []*responses.ViolationsDetail{
					{
						Name:    "field1",
						Message: "Invalid value",
					},
				},
			},
			expectedCode:    "INVALID_FORMAT",
			expectedMessage: "Test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.message, tt.validationErrs)
			valErr, ok := err.(*responses.ValidationError)

			if !ok {
				t.Errorf("Expected ValidationError, got %T", err)
				return
			}

			if valErr.Code != tt.expectedCode {
				t.Errorf("Expected code %s, got %s", tt.expectedCode, valErr.Code)
			}

			if valErr.Message != tt.expectedMessage {
				t.Errorf("Expected message %s, got %s", tt.expectedMessage, valErr.Message)
			}

			if tt.validationErrs != nil {
				if len(valErr.Parameters) != len(tt.validationErrs.Parameters) {
					t.Errorf("Expected %d parameters, got %d",
						len(tt.validationErrs.Parameters),
						len(valErr.Parameters))
				}
				if len(valErr.Violations) != len(tt.validationErrs.Violations) {
					t.Errorf("Expected %d violations, got %d",
						len(tt.validationErrs.Violations),
						len(valErr.Violations))
				}
			}
		})
	}
}

func TestGetErrorMessage(t *testing.T) {
	tests := []struct {
		name            string
		fieldError      validator.FieldError
		expectedMessage string
	}{
		{
			name:            "required field",
			fieldError:      mockFieldError{tag: "required"},
			expectedMessage: "This field is required",
		},
		{
			name:            "email validation",
			fieldError:      mockFieldError{tag: "email"},
			expectedMessage: "Invalid email format",
		},
		{
			name:            "min validation",
			fieldError:      mockFieldError{tag: "min", param: "5"},
			expectedMessage: "Should be at least 5",
		},
		{
			name:            "max validation",
			fieldError:      mockFieldError{tag: "max", param: "10"},
			expectedMessage: "Should be at most 10",
		},
		{
			name:            "ascii validation",
			fieldError:      mockFieldError{tag: "ascii"},
			expectedMessage: "Should contain only ASCII characters",
		},
		{
			name:            "gt validation with 0",
			fieldError:      mockFieldError{tag: "gt", param: "0"},
			expectedMessage: "Should be a positive number",
		},
		{
			name:            "gt validation with other value",
			fieldError:      mockFieldError{tag: "gt", param: "5"},
			expectedMessage: "Should be greater than 5",
		},
		{
			name:            "unknown validation",
			fieldError:      mockFieldError{tag: "unknown"},
			expectedMessage: "Failed validation: unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := GetErrorMessage(tt.fieldError)
			if message != tt.expectedMessage {
				t.Errorf("Expected message %q, got %q", tt.expectedMessage, message)
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name: "is validation error",
			err: &responses.ValidationError{
				Code:    "INVALID_FORMAT",
				Message: "Test error",
			},
			expected: true,
		},
		{
			name:     "not validation error",
			err:      fmt.Errorf("regular error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, result := responses.IsValidationError(tt.err)
			if result != tt.expected {
				t.Errorf("Expected IsValidationError to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHandleValidationErrors(t *testing.T) {
	tests := []struct {
		name           string
		inputErr       error
		expectedError  bool
		expectedFields []string
	}{
		{
			name:           "nil error",
			inputErr:       nil,
			expectedError:  false,
			expectedFields: nil,
		},
		{
			name:           "non-validation error",
			inputErr:       fmt.Errorf("regular error"),
			expectedError:  true,
			expectedFields: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HandleValidationErrors(tt.inputErr)

			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if tt.expectedError && err == nil {
				t.Error("Expected error, got nil")
				return
			}

			if valErr, ok := err.(*responses.ValidationError); ok {
				if tt.expectedFields != nil {
					if len(valErr.Parameters) != len(tt.expectedFields) {
						t.Errorf("Expected %d fields, got %d",
							len(tt.expectedFields),
							len(valErr.Parameters))
					}

					for i, field := range tt.expectedFields {
						if valErr.Parameters[i] != field {
							t.Errorf("Expected field %s at position %d, got %s",
								field, i, valErr.Parameters[i])
						}
					}
				}
			}
		})
	}
}
