package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	rule "github.com/lozhkindm/otus-go-hw/hw09_struct_validator/extractors"
)

var (
	errNotStructure = errors.New("the given value is not a structure")
	extractors      = make(map[string]rule.Extractor)
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	for i, vErr := range v {
		builder.WriteString(vErr.Err.Error())
		if i != len(v)-1 {
			builder.WriteRune('\n')
		}
	}
	return builder.String()
}

func Validate(v interface{}, prefix string) error {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return errNotStructure
	}
	iv := reflect.ValueOf(v)
	vErrors := make(ValidationErrors, 0, iv.NumField())
	var err error
	for i := 0; i < iv.Type().NumField(); i++ {
		vErrors, err = validateStructField(iv.Type().Field(i), iv.Field(i), vErrors, prefix)
		if err != nil {
			var ve ValidationErrors
			if errors.As(err, &ve) {
				vErrors = append(vErrors, ve...)
			} else {
				return err
			}
		}
	}
	if len(vErrors) > 0 {
		return vErrors
	}
	return nil
}

func validateStructField(
	field reflect.StructField,
	value reflect.Value,
	vErrors ValidationErrors,
	prefix string,
) (ValidationErrors, error) {
	if !field.IsExported() {
		return vErrors, nil
	}
	if isStructToValidate(field) {
		prefix = fmt.Sprintf("%s%s->", prefix, field.Name)
		return vErrors, Validate(value.Interface(), prefix)
	}
	ext := getRulesExtractor(field.Type.Kind(), field.Type)
	if ext == nil {
		return vErrors, nil
	}
	rules, err := ext.Extract(field, value)
	if err != nil {
		return vErrors, err
	}
	for _, r := range rules {
		valid, err := r.Validate()
		if err != nil {
			return vErrors, err
		}
		if !valid {
			vErrors = append(vErrors, ValidationError{
				Field: field.Name,
				Err:   r.ValidationError(prefix + field.Name),
			})
		}
	}
	return vErrors, nil
}

func isStructToValidate(field reflect.StructField) bool {
	if field.Type.Kind() != reflect.Struct {
		return false
	}
	tag, ok := field.Tag.Lookup("validate")
	if !ok {
		return false
	}
	return tag == "nested"
}

func getRulesExtractor(kind reflect.Kind, fieldType reflect.Type) rule.Extractor {
	switch kind { // nolint: exhaustive
	case reflect.String:
		if _, ok := extractors["string"]; !ok {
			extractors["string"] = rule.NewStringRulesExtractor("validate", "|", ":", ",")
		}
		return extractors["string"]
	case reflect.Int:
		if _, ok := extractors["int"]; !ok {
			extractors["int"] = rule.NewIntRulesExtractor("validate", "|", ":", ",")
		}
		return extractors["int"]
	case reflect.Slice:
		return getSliceRulesExtractor(fieldType)
	default:
		return nil
	}
}

func getSliceRulesExtractor(fieldType reflect.Type) rule.Extractor {
	switch fieldType.String() {
	case "[]string":
		if _, ok := extractors["string"]; !ok {
			extractors["string"] = rule.NewStringRulesExtractor("validate", "|", ":", ",")
		}
		return extractors["string"]
	case "[]int":
		if _, ok := extractors["int"]; !ok {
			extractors["int"] = rule.NewIntRulesExtractor("validate", "|", ":", ",")
		}
		return extractors["int"]
	}
	return nil
}
