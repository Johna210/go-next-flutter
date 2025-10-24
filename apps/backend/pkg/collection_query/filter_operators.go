package collectionquery

type FilterOperators string

const (
	EqualTo              FilterOperators = "="
	Between              FilterOperators = "BETWEEN"
	LessThan             FilterOperators = "<"
	LessThanOrEqualTo    FilterOperators = "<="
	GreaterThan          FilterOperators = ">"
	GreaterThanOrEqualTo FilterOperators = ">="
	In                   FilterOperators = "IN"
	NotIn                FilterOperators = "NotIn"
	Any                  FilterOperators = "ANY"
	NotNull              FilterOperators = "NotNull"
	IsNotNull            FilterOperators = "IsNotNull"
	IsNull               FilterOperators = "IsNull"
	NotEqualTo           FilterOperators = "!="
	Like                 FilterOperators = "LIKE"
	ILike                FilterOperators = "ILIKE"
	NotEqual             FilterOperators = "NotEqual"
	All                  FilterOperators = "All"
	ArrayFilter          FilterOperators = "ArrayFilter"
	ArrayContains        FilterOperators = "ArrayContains"
)

type FilterSeparators string

const (
	WhereEqual FilterSeparators = "_:"
	WhereAND   FilterSeparators = "_|"
	WhereOR    FilterSeparators = "_,"
	OrderBy    FilterSeparators = ","
	OrderItem  FilterSeparators = ":"
)
