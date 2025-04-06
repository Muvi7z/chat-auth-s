package validate

import (
	"encoding/json"
	"errors"
)

type ValidationError struct {
	Message []string `json:"error_message"`
}

func NewValidationErrors(message ...string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}

func (v *ValidationError) addError(message string) {
	v.Message = append(v.Message, message)
}

func (v *ValidationError) Error() string {
	data, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}

	return string(data)
}

func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}
