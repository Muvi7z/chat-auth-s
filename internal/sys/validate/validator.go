package validate

import "context"

type Condition func(ctx context.Context) error

func Validate(ctx context.Context, conditions ...Condition) error {
	ve := NewValidationErrors()

	for _, condition := range conditions {
		err := condition(ctx)
		if err != nil {
			if IsValidationError(err) {
				ve.addError(err.Error())
				continue
			}

			return err
		}
	}

	if ve.Message == nil {
		return nil
	}

	return ve
}
