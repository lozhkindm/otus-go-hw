package rules

import (
	"errors"
	"fmt"
)

type Max struct {
	value   int64
	max     int64
	verbose bool
}

func NewRuleMax(value, max int64, verbose bool) Max {
	return Max{
		value:   value,
		max:     max,
		verbose: verbose,
	}
}

func (m Max) Validate() (bool, error) {
	return m.value <= m.max, nil
}

func (m Max) ValidationError(field string) error {
	if m.verbose {
		return errors.New(fmt.Sprintf("the field %q (%d) may not be greater than %d", field, m.value, m.max))
	}
	return errors.New(fmt.Sprintf("the field %q may not be greater than %d", field, m.max))
}
