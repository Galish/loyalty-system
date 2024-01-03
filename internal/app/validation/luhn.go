package validation

import (
	"errors"

	"github.com/ShiraazMoollatjie/goluhn"
)

var ErrInvalidOrderNumber = errors.New("invalid order number value")

func LuhnValidate(num string) error {
	if num == "" {
		return ErrInvalidOrderNumber
	}

	if err := goluhn.Validate(num); err != nil {
		return ErrInvalidOrderNumber
	}

	return nil
}
