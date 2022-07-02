package extractors

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/lozhkindm/otus-go-hw/hw09_struct_validator/rules"
)

const (
	ruleMin   = "min"
	ruleMax   = "max"
	ruleRange = "in"
)

var (
	ErrInInvalidOptions = errors.New(`rule "in" must have integer min and max values`)
	ErrInvalidMinValue  = errors.New(`rule "min" must have integer value`)
	ErrInvalidMaxValue  = errors.New(`rule "max" must have integer value`)
)

type IntRules struct {
	tagName        string
	tagSeparator   string
	ruleSeparator  string
	valueSeparator string
}

func NewIntRulesExtractor(tagName, tagSeparator, ruleSeparator, valueSeparator string) IntRules {
	return IntRules{
		tagName:        tagName,
		tagSeparator:   tagSeparator,
		ruleSeparator:  ruleSeparator,
		valueSeparator: valueSeparator,
	}
}

func (i IntRules) Extract(field reflect.StructField, value reflect.Value) ([]rules.Rule, error) {
	tag, ok := field.Tag.Lookup(i.tagName)
	if !ok {
		return nil, nil
	}
	rulesStr := strings.Split(tag, i.tagSeparator)
	extracted := make([]rules.Rule, 0, len(rulesStr))
	for _, ruleStr := range rulesStr {
		ruleParts := strings.Split(ruleStr, i.ruleSeparator)
		ruleName, ruleVal := ruleParts[0], ruleParts[1]
		switch ruleName {
		case ruleMin:
			min, err := strconv.Atoi(ruleVal)
			if err != nil {
				return nil, ErrInvalidMinValue
			}
			extracted = addRuleMin(extracted, value, min)
		case ruleMax:
			max, err := strconv.Atoi(ruleVal)
			if err != nil {
				return nil, ErrInvalidMaxValue
			}
			extracted = addRuleMax(extracted, value, max)
		case ruleRange:
			if ruleVal == "" {
				return nil, ErrInNoOptions
			}
			valsString := strings.Split(ruleVal, i.valueSeparator)
			vals := make([]int64, 0, len(valsString))
			for _, valStr := range valsString {
				val, err := strconv.Atoi(valStr)
				if err != nil {
					return nil, ErrInInvalidOptions
				}
				vals = append(vals, int64(val))
			}
			extracted = addRuleRange(extracted, value, vals)
		}
	}
	return extracted, nil
}

func addRuleMin(extracted []rules.Rule, value reflect.Value, min int) []rules.Rule {
	if value.Type().String() == "[]int" {
		for j := 0; j < value.Len(); j++ {
			extracted = append(extracted, rules.NewRuleMin(value.Index(j).Int(), int64(min), true))
		}
	} else {
		extracted = append(extracted, rules.NewRuleMin(value.Int(), int64(min), false))
	}
	return extracted
}

func addRuleMax(extracted []rules.Rule, value reflect.Value, max int) []rules.Rule {
	if value.Type().String() == "[]int" {
		for j := 0; j < value.Len(); j++ {
			extracted = append(extracted, rules.NewRuleMax(value.Index(j).Int(), int64(max), true))
		}
	} else {
		extracted = append(extracted, rules.NewRuleMax(value.Int(), int64(max), false))
	}
	return extracted
}

func addRuleRange(extracted []rules.Rule, value reflect.Value, vals []int64) []rules.Rule {
	if value.Type().String() == "[]int" {
		for j := 0; j < value.Len(); j++ {
			extracted = append(extracted, rules.NewRuleRange(value.Index(j).Int(), vals, true))
		}
	} else {
		extracted = append(extracted, rules.NewRuleRange(value.Int(), vals, false))
	}
	return extracted
}
