package responses

func NewValidationError(message string, parameters []string, violations []*ViolationsDetail) *ValidationError {
	return &ValidationError{
		Code:       "INVALID_FORMAT",
		Message:    message,
		Parameters: parameters,
		Violations: violations,
	}
}

func NewNotFoundError(message string) *ValidationError {
	return &ValidationError{
		Code:    "NOT_FOUND",
		Message: message,
	}
}

func NewConflictError(message string) *ValidationError {
	return &ValidationError{
		Code:    "CONFLICT",
		Message: message,
	}
}

func NewInternalServerError(message string) *ValidationError {
	return &ValidationError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: message,
	}
}
