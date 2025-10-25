package collectionquery

import "github.com/go-playground/validator/v10"

var ValidateFilterOperators = map[FilterOperators]bool{
	EqualTo:              true,
	Between:              true,
	LessThan:             true,
	LessThanOrEqualTo:    true,
	GreaterThan:          true,
	GreaterThanOrEqualTo: true,
	In:                   true,
	NotIn:                true,
	Any:                  true,
	NotNull:              true,
	IsNotNull:            true,
	IsNull:               true,
	NotEqualTo:           true,
	Like:                 true,
	ILike:                true,
	NotEqual:             true,
	All:                  true,
	ArrayFilter:          true,
	ArrayContains:        true,
}

func RegisterFilterOperatorValidator(v *validator.Validate) {
	_ = v.RegisterValidation("valid_filter_operators", func(fl validator.FieldLevel) bool {
		op := FilterOperators(fl.Field().String())
		return ValidateFilterOperators[op]
	})
}
