package collectionquery

import "github.com/go-playground/validator/v10"

var ValidateFilterOperator = map[FilterOperator]bool{
	OpEq:             true,
	OpGt:             true,
	OpGte:            true,
	OpLt:             true,
	OpLte:            true,
	OpNotEqAngle:     true,
	OpNotEq:          true,
	OpLike:           true,
	OpILike:          true,
	OpRegex:          true,
	OpIRegex:         true,
	OpIn:             true,
	OpIs:             true,
	OpIsDistinctFrom: true,
	OpTsQuery:        true,
	OpContains:       true,
	OpIsContainedBy:  true,
	OpOverlaps:       true,
	OpNotExtendRight: true,
	OpNotExtendLeft:  true,
	OpAdjacent:       true,
	OpNot:            true,
	OpOr:             true,
	OpAnd:            true,
	OpAll:            true,
	OpAny:            true,
	OpBetween:        true,
}

func RegisterFilterOperatorValidator(v *validator.Validate) {
	_ = v.RegisterValidation("valid_filter_operator", func(fl validator.FieldLevel) bool {
		op := FilterOperator(fl.Field().String())
		return ValidateFilterOperator[op]
	})
}
