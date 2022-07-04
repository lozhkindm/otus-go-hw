package rules

import (
	"errors"
	"fmt"
)

type Min struct {
	value   int64
	min     int64
	verbose bool
}

func NewRuleMin(value, min int64, verbose bool) Min {
	return Min{
		value:   value,
		min:     min,
		verbose: verbose,
	}
}

func (m Min) Validate() (bool, error) {
	return m.value >= m.min, nil
}

func (m Min) ValidationError(field string) error {
	if m.verbose {
		return errors.New(fmt.Sprintf("the field %q (%d) must be at least %d", field, m.value, m.min))
	}
	return errors.New(fmt.Sprintf("the field %q must be at least %d", field, m.min))
}
