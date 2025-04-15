package responses

type ValidationError struct {
	Code       string              `json:"code"`
	Message    string              `json:"message"`
	Parameters []string            `json:"parameters"`
	Violations []*ViolationsDetail `json:"violations"`
}

func (e ValidationError) Error() string {
	return e.Message
}

type ValidationErrors struct {
	Parameters []string
	Violations []*ViolationsDetail
}

type ViolationsDetail struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func IsValidationError(err error) (*ValidationError, bool) {
	if ve, ok := err.(ValidationError); ok {
		return &ve, true
	}

	if ve, ok := err.(*ValidationError); ok {
		return ve, true
	}

	return nil, false
}
