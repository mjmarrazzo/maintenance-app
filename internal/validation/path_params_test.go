package validation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mjmarrazzo/maintenance-app/internal/responses"
)

type PathParamTestStruct struct {
	ID       int     `param:"id" validate:"required,gt=0"`
	UserID   uint    `param:"user_id" validate:"required"`
	Rating   float64 `param:"rating" validate:"gte=0,lte=5"`
	IsActive bool    `param:"is_active"`
	Name     string  `param:"name" validate:"required"`
	Optional string  `param:"optional"`
}

func TestBindPathParams(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		params       map[string]string
		expectError  bool
		errorCode    string
		errorMessage string
		validate     func(t *testing.T, params *PathParamTestStruct)
	}{
		{
			name: "valid parameters",
			path: "/users/:id/:user_id/:rating/:is_active/:name",
			params: map[string]string{
				"id":        "123",
				"user_id":   "456",
				"rating":    "4.5",
				"is_active": "true",
				"name":      "test",
			},
			expectError: false,
			validate: func(t *testing.T, params *PathParamTestStruct) {
				if params.ID != 123 {
					t.Errorf("Expected ID 123, got %d", params.ID)
				}
				if params.UserID != 456 {
					t.Errorf("Expected UserID 456, got %d", params.UserID)
				}
				if params.Rating != 4.5 {
					t.Errorf("Expected Rating 4.5, got %f", params.Rating)
				}
				if !params.IsActive {
					t.Error("Expected IsActive true, got false")
				}
				if params.Name != "test" {
					t.Errorf("Expected Name 'test', got '%s'", params.Name)
				}
			},
		},
		{
			name: "invalid integer",
			path: "/users/:id",
			params: map[string]string{
				"id": "not-a-number",
			},
			expectError:  true,
			errorCode:    "INVALID_FORMAT",
			errorMessage: "Validation failed",
		},
		{
			name: "invalid uint",
			path: "/users/:user_id",
			params: map[string]string{
				"user_id": "-1",
			},
			expectError:  true,
			errorCode:    "INVALID_FORMAT",
			errorMessage: "Validation failed",
		},
		{
			name: "invalid float",
			path: "/users/:rating",
			params: map[string]string{
				"rating": "invalid-float",
			},
			expectError:  true,
			errorCode:    "INVALID_FORMAT",
			errorMessage: "Validation failed",
		},
		{
			name: "invalid boolean",
			path: "/users/:is_active",
			params: map[string]string{
				"is_active": "not-a-bool",
			},
			expectError:  true,
			errorCode:    "INVALID_FORMAT",
			errorMessage: "Validation failed",
		},
		{
			name: "missing required field",
			path: "/users",
			params: map[string]string{
				"optional": "value",
			},
			expectError:  true,
			errorCode:    "INVALID_FORMAT",
			errorMessage: "Validation failed",
		},
		{
			name: "validation error - ID less than 1",
			path: "/users/:id",
			params: map[string]string{
				"id": "0",
			},
			expectError:  true,
			errorCode:    "INVALID_FORMAT",
			errorMessage: "Validation failed",
		},
		{
			name: "validation error - rating out of range",
			path: "/users/:id/:rating",
			params: map[string]string{
				"id":     "1",
				"rating": "6.0",
			},
			expectError:  true,
			errorCode:    "INVALID_FORMAT",
			errorMessage: "Validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set path parameters
			names := make([]string, 0)
			values := make([]string, 0)
			for name, value := range tt.params {
				names = append(names, name)
				values = append(values, value)
			}
			c.SetParamNames(names...)
			c.SetParamValues(values...)

			// Test target
			var params PathParamTestStruct
			err := BindPathParams(c, &params)

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

				if tt.validate != nil {
					tt.validate(t, &params)
				}
			}
		})
	}
}
