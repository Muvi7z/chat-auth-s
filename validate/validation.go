package validate

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/internal/sys/validate"
)

func ValidateId(id int64) validate.Condition {
	return func(ctx context.Context) error {
		if id <= 0 {
			return validate.NewValidationErrors("id must be a positive number")
		}

		return nil
	}
}

func OtherValidateID(id int64) validate.Condition {
	return func(ctx context.Context) error {
		if id <= 100 {
			return validate.NewValidationErrors("id must be greater than 100")
		}

		return nil
	}
}
