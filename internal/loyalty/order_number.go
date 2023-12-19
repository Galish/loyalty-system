package loyalty

import "github.com/ShiraazMoollatjie/goluhn"

type OrderNumber string

func (num OrderNumber) isValid() bool {
	if num.String() == "" {
		return false
	}

	if err := goluhn.Validate(num.String()); err != nil {
		return false
	}

	return true
}

func (num OrderNumber) String() string {
	return string(num)
}
