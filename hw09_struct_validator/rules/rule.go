package rules

type Rule interface {
	Validate() (bool, error)
	ValidationError(field string) error
}
