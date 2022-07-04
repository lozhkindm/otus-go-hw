package extractors

import (
	"reflect"

	"github.com/lozhkindm/otus-go-hw/hw09_struct_validator/rules"
)

type Extractor interface {
	Extract(field reflect.StructField, value reflect.Value) ([]rules.Rule, error)
}
