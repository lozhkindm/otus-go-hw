package rules

import (
	"errors"
	"fmt"
	"strings"
)

type In struct {
	string  string
	values  []string
	verbose bool
}

func NewRuleIn(string string, values []string, verbose bool) In {
	return In{
		string:  string,
		values:  values,
		verbose: verbose,
	}
}

func (i In) Validate() (bool, error) {
	in := false
	for _, v := range i.values {
		if i.string == v {
			in = true
			break
		}
	}
	return in, nil
}

func (i In) ValidationError(field string) error {
	if i.verbose {
		return errors.New(fmt.Sprintf("the field %q (%s) must have one of the following values: %s", field, i.string, strings.Join(i.values, ", ")))
	}
	return errors.New(fmt.Sprintf("the field %q must have one of the following values: %s", field, strings.Join(i.values, ", ")))
}
