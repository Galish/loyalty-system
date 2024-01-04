package validation

import (
	"github.com/ShiraazMoollatjie/goluhn"
)

func IsValidLuhn(num string) bool {
	if num == "" {
		return false
	}

	if err := goluhn.Validate(num); err != nil {
		return false
	}

	return true
}
