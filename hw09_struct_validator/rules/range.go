package rules

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Range struct {
	int     int64
	values  []int64
	verbose bool
}

func NewRuleRange(int int64, values []int64, verbose bool) Range {
	return Range{
		int:     int,
		values:  values,
		verbose: verbose,
	}
}

func (r Range) Validate() (bool, error) {
	in := false
	for _, v := range r.values {
		if r.int == v {
			in = true
			break
		}
	}
	return in, nil
}

func (r Range) ValidationError(field string) error {
	if r.verbose {
		return errors.New(fmt.Sprintf("the field %q (%d) must have one of the following values: %s", field, r.int, r.getValuesString()))
	}
	return errors.New(fmt.Sprintf("the field %q must have one of the following values: %s", field, r.getValuesString()))
}

func (r Range) getValuesString() string {
	var builder strings.Builder
	for i, v := range r.values {
		builder.WriteString(strconv.Itoa(int(v)))
		if i != len(r.values)-1 {
			builder.WriteString(", ")
		}
	}
	return builder.String()
}
