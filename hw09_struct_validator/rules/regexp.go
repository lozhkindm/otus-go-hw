package rules

import (
	"errors"
	"fmt"
	"regexp"
)

type Regexp struct {
	string  string
	expr    string
	verbose bool
}

func NewRuleRegexp(string, expr string, verbose bool) Regexp {
	return Regexp{
		string:  string,
		expr:    expr,
		verbose: verbose,
	}
}

func (r Regexp) Validate() (bool, error) {
	reg, err := regexp.Compile(r.expr)
	if err != nil {
		return false, err
	}
	return reg.Match([]byte(r.string)), nil
}

func (r Regexp) ValidationError(field string) error {
	if r.verbose {
		return errors.New(fmt.Sprintf("the field %q (%s) format is invalid", field, r.string))
	}
	return errors.New(fmt.Sprintf("the field %q format is invalid", field))
}
