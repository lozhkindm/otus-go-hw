package rules

import (
	"errors"
	"fmt"
)

type Len struct {
	string  string
	length  int
	verbose bool
}

func NewRuleLen(string string, length int, verbose bool) Len {
	return Len{
		string:  string,
		length:  length,
		verbose: verbose,
	}
}

func (l Len) Validate() (bool, error) {
	return len(l.string) == l.length, nil
}

func (l Len) ValidationError(field string) error {
	if l.verbose {
		return errors.New(fmt.Sprintf("the field %q (%s) must be %d characters", field, l.string, l.length))
	}
	return errors.New(fmt.Sprintf("the field %q must be %d characters", field, l.length))
}
