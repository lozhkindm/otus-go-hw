package extractors

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/lozhkindm/otus-go-hw/hw09_struct_validator/rules"
)

const (
	ruleLen    = "len"
	ruleRegexp = "regexp"
	ruleIn     = "in"
)

var (
	ErrInNoOptions     = errors.New(`no options were specified for "in" rule`)
	ErrInvalidLenValue = errors.New(`rule "len" must have integer value`)
)

type StringRules struct {
	tagName        string
	tagSeparator   string
	ruleSeparator  string
	valueSeparator string
}

func NewStringRulesExtractor(tagName, tagSeparator, ruleSeparator, valueSeparator string) StringRules {
	return StringRules{
		tagName:        tagName,
		tagSeparator:   tagSeparator,
		ruleSeparator:  ruleSeparator,
		valueSeparator: valueSeparator,
	}
}

func (s StringRules) Extract(field reflect.StructField, value reflect.Value) ([]rules.Rule, error) {
	tag, ok := field.Tag.Lookup(s.tagName)
	if !ok {
		return nil, nil
	}
	rulesStr := strings.Split(tag, s.tagSeparator)
	extracted := make([]rules.Rule, 0, len(rulesStr))
	for _, ruleStr := range rulesStr {
		ruleParts := strings.Split(ruleStr, s.ruleSeparator)
		ruleName, ruleVal := ruleParts[0], ruleParts[1]
		switch ruleName {
		case ruleLen:
			length, err := strconv.Atoi(ruleVal)
			if err != nil {
				return nil, ErrInvalidLenValue
			}
			extracted = addRuleLen(extracted, value, length)
		case ruleRegexp:
			extracted = addRuleRegexp(extracted, value, ruleVal)
		case ruleIn:
			if ruleVal == "" {
				return nil, ErrInNoOptions
			}
			vals := strings.Split(ruleVal, s.valueSeparator)
			extracted = addRuleIn(extracted, value, vals)
		}
	}
	return extracted, nil
}

func addRuleLen(extracted []rules.Rule, value reflect.Value, length int) []rules.Rule {
	if value.Type().String() == "[]string" {
		for j := 0; j < value.Len(); j++ {
			extracted = append(extracted, rules.NewRuleLen(value.Index(j).String(), length, true))
		}
	} else {
		extracted = append(extracted, rules.NewRuleLen(value.String(), length, false))
	}
	return extracted
}

func addRuleRegexp(extracted []rules.Rule, value reflect.Value, ruleVal string) []rules.Rule {
	if value.Type().String() == "[]string" {
		for j := 0; j < value.Len(); j++ {
			extracted = append(extracted, rules.NewRuleRegexp(value.Index(j).String(), ruleVal, true))
		}
	} else {
		extracted = append(extracted, rules.NewRuleRegexp(value.String(), ruleVal, false))
	}
	return extracted
}

func addRuleIn(extracted []rules.Rule, value reflect.Value, ruleValues []string) []rules.Rule {
	if value.Type().String() == "[]string" {
		for j := 0; j < value.Len(); j++ {
			extracted = append(extracted, rules.NewRuleIn(value.Index(j).String(), ruleValues, true))
		}
	} else {
		extracted = append(extracted, rules.NewRuleIn(value.String(), ruleValues, false))
	}
	return extracted
}
