package validation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mjmarrazzo/maintenance-app/internal/responses"
)

type TestStruct struct {
	Name     string `json:"name" validate:"required"`
	Age      int    `json:"age" validate:"gte=0"`
	Email    string `json:"email" validate:"required,email"`
	Optional string `json:"optional"`
}

func TestBindBody(t *testing.T) {
	tests := []struct {
		name         string
		payload      string
		expectError  bool
		errorCode    string
		errorMessage string
		validateBody func(t *testing.T, body *TestStruct)
	}{
		{
			name:        "valid request",
			payload:     `{"name": "John", "age": 25, "email": "john@example.com"}`,
			expectError: false,
			validateBody: func(t *testing.T, body *TestStruct) {
				if body.Name != "John" {
					t.Errorf("Expected name 'John', got '%s'", body.Name)
				}
				if body.Age != 25 {
					t.Errorf("Expected age 25, got %d", body.Age)
				}
				if body.Email != "john@example.com" {
					t.Errorf("Expected email 'john@example.com', got '%s'", body.Email)
				}
			},
		},
		{
			name:         "missing required field",
			payload:      `{"age": 25, "email": "john@example.com"}`,
			expectError:  true,
			errorCode:    "INVALID_FORMAT",
			errorMessage: "Validation failed",
		},
		{
			name:         "invalid email format",
			payload:      `{"name": "John", "age": 25, "email": "invalid-email"}`,
			expectError:  true,
			errorCode:    "INVALID_FORMAT",
			errorMessage: "Validation failed",
		},
		{
			name:         "negative age",
			payload:      `{"name": "John", "age": -1, "email": "john@example.com"}`,
			expectError:  true,
			errorCode:    "INVALID_FORMAT",
			errorMessage: "Validation failed",
		},
		{
			name:         "extra field",
			payload:      `{"name": "John", "age": 25, "email": "john@example.com", "unknown": "value"}`,
			expectError:  true,
			errorCode:    "INVALID_FORMAT",
			errorMessage: "Validation failed",
		},
		{
			name:         "invalid json",
			payload:      `{"name": "John", "age": "not a number", "email": "john@example.com"}`,
			expectError:  true,
			errorCode:    "INVALID_FORMAT",
			errorMessage: "Validation failed",
		},
		// malformed json
		{
			name:         "malformed json",
			payload:      `{"name": "John", "age": "not a number" "email": "john@example.com"}`,
			expectError:  true,
			errorCode:    "INVALID_FORMAT",
			errorMessage: "Invalid JSON syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Test target
			var body TestStruct
			err := BindBody(c, &body)

			// Assertions
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}

				var valErr *responses.ValidationError
				if ve, ok := responses.IsValidationError(err); !ok {
					t.Errorf("Expected ValidationError, got %T", err)
					return
				} else {
					valErr = ve
				}

				if valErr.Code != tt.errorCode {
					t.Errorf("Expected error code %s, got %s", tt.errorCode, valErr.Code)
				}
				if valErr.Message != tt.errorMessage {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMessage, valErr.Message)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
					return
				}

				if tt.validateBody != nil {
					tt.validateBody(t, &body)
				}
			}
		})
	}
}

func TestValidateNoExtraFields(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		expectError bool
	}{
		{
			name:        "no extra fields",
			payload:     `{"name": "John", "age": 25, "email": "john@example.com"}`,
			expectError: false,
		},
		{
			name:        "with extra field",
			payload:     `{"name": "John", "age": 25, "email": "john@example.com", "extra": "field"}`,
			expectError: true,
		},
		{
			name:        "with optional field",
			payload:     `{"name": "John", "age": 25, "email": "john@example.com", "optional": "value"}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Test target
			var ts TestStruct
			err := validateNoExtraFields(c, &ts)

			// Assertions
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestHandleBindError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedCode   int
		expectedFields []string
	}{
		{
			name: "unmarshal type error",
			err: &json.UnmarshalTypeError{
				Field: "age",
				Type:  reflect.TypeOf(0),
				Value: "string",
			},
			expectedCode:   http.StatusBadRequest,
			expectedFields: []string{"age"},
		},
		{
			name: "validation error",
			err: &responses.ValidationError{
				Code:       "INVALID_FORMAT",
				Message:    "Test error",
				Parameters: []string{"name"},
				Violations: []*responses.ViolationsDetail{
					{Name: "name", Message: "Required field"},
				},
			},
			expectedCode:   http.StatusBadRequest,
			expectedFields: []string{"name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Test target
			err := HandleBindError(c, tt.err, reflect.TypeOf(TestStruct{}))

			// Assertions
			if err == nil {
				t.Error("Expected error but got none")
				return
			}

			var valErr *responses.ValidationError
			if ve, ok := responses.IsValidationError(err); !ok {
				t.Errorf("Expected ValidationError, got %T", err)
				return
			} else {
				valErr = ve
			}

			if len(tt.expectedFields) > 0 {
				for _, expectedField := range tt.expectedFields {
					found := false
					for _, param := range valErr.Parameters {
						if param == expectedField {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected field %s in parameters, but not found", expectedField)
					}
				}
			}
		})
	}
}
