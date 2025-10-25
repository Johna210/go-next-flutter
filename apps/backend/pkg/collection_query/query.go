package collectionquery

// SortDirection defines the allowed sorting direction
type SortDirection string

const (
	Ascending  SortDirection = "ASC"
	Descending SortDirection = "DESC"
)

// NullsOrder defines the allowed ordering for NULL values
type NullsOrder string

const (
	NullsFirst NullsOrder = "NULLS_FIRST"
	NullsLast  NullsOrder = "NULLS_LAST"
)

type FilterOperator string

const (
	OpEq             FilterOperator = "="
	OpGt             FilterOperator = ">"
	OpGte            FilterOperator = ">="
	OpLt             FilterOperator = "<"
	OpLte            FilterOperator = "<="
	OpNotEqAngle     FilterOperator = "<>"
	OpNotEq          FilterOperator = "!="
	OpLike           FilterOperator = "LIKE"
	OpILike          FilterOperator = "ILIKE"
	OpRegex          FilterOperator = "~"
	OpIRegex         FilterOperator = "~*"
	OpIn             FilterOperator = "IN"
	OpIs             FilterOperator = "IS"
	OpIsDistinctFrom FilterOperator = "IS DISTINCT FROM"
	OpTsQuery        FilterOperator = "@@"
	OpContains       FilterOperator = "contains"
	OpIsContainedBy  FilterOperator = "<@"
	OpOverlaps       FilterOperator = "&&"
	OpNotExtendRight FilterOperator = "&<"
	OpNotExtendLeft  FilterOperator = "&>"
	OpAdjacent       FilterOperator = "-|-"
	OpNot            FilterOperator = "NOT"
	OpOr             FilterOperator = "OR"
	OpAnd            FilterOperator = "AND"
	OpAll            FilterOperator = "ALL"
	OpAny            FilterOperator = "ANY"
	OpBetween        FilterOperator = "BETWEEN"
)

// CollectionQuery defines the parameters for fetching, filtering and shaping
// a collection of results from a data source.
type CollectionQuery struct {
	// Select specifies which columns to return. If empty all columns are returned.
	// Example: []string{"id", "name", "email"}
	Select []string `json:"select,omitempty" validate:"omitempty,dive,required"`

	// Where defines the filtering conditions. It's a 2D slice to handle
	// complex AND/OR logic: outer slices are joined by AND, inner slices by OR.
	// Example: [[{Column: "status", Op: "=", Val: "active"}], [{Column: "age", Op: ">", Val: 21}]]
	// This translates to: (status = 'active') AND (age > 21)
	Where [][]Where `json:"where,omitempty" validate:"omitempty,dive,dive"`

	// Take sets the maximum number of records to return (limit).
	// Use a pointer to distinguish between zero value and not set.
	Take *int `json:"take,omitempty" validate:"omitempty,min=0"`

	// Skip sets the number of records to skip (offset).
	Skip *int `json:"skip,omitempty" validate:"omitempty,min=0"`

	// OrderBy specifies the sorting order for the results.
	OrderBy []Order `json:"order_by,omitempty" validate:"omitempty,dive"`

	// Includes specifies related entities to include in the result.
	// The names should correspond to predefined relations in the data model.
	// Example: []string{"profile", "orders"}
	Includes []string `json:"includes,omitempty" validate:"omitempty,dive,required"`

	// IncludeAndSelect allows specifying which columns to select from included relations.
	// Example: []IncludeSelect{{Name: "profile", Select: []string{"first_name", "last_name"}}}
	IncludeAndSelect []IncludeSelect `json:"include_and_select,omitempty" validate:"omitempty,dive"`

	// LeftJoinAndMapOne provides a way to define custom joins.
	// The 'any' type is used here as a placeholder for a more specific struct
	// that you would define to match your application's needs.
	LeftJoinAndMapOne []any `json:"left_join_and_map_one,omitempty" validate:"omitempty,dive"`

	// GroupBy specifies the columns to group the results by. typically for aggregate functions.
	GroupBy []string `json:"group_by,omitempty" validate:"omitempty,dive,required"`

	// Having adds filtering conditions on grouped rows (used with GroupBy).
	// It follows the same AND/OR logic as the Where field.
	Having [][]Where `json:"having,omitempty" validate:"omitempty,dive,dive"`

	// Count, if true, makes the query return only the count of matching records
	// records instead of the data itself.
	Count *bool `json:"count,omitempty" validate:"omitempty"`
}

// Order specifies the column and direction for sorting a query result.
// It includes validation tags for ensuring data integrity.
type Order struct {
	// Column is the name of the database column to sort by.
	// It is required field.
	Column string `json:"column" validate:"required"`

	// Direction specifies the sort direction, either ascending (ASC) or descending (DESC).
	// It is an optional field.
	Direction *SortDirection `json:"direction,omitempty" validate:"omitempty,oneof=ASC DESC"`

	// Nulls specifies how NULL values should be ordered, either first (NULLS_FIRST) or last (NULLS_LAST).
	// It is an optional field.
	Nulls *NullsOrder `json:"nulls,omitempty" validate:"omitempty,oneof=NULLS_FIRST NULLS_LAST"`
}

// Where defines a single condition for a filter, consisting of a
// column, an operator, and a value to compare against.
type Where struct {
	// Column is the name of the database column to apply the filter to.
	// This field is required.
	Column string `json:"column" validate:"required"`

	// Operator is the comparison operator to use for the filter.
	// This field is required.
	Operator FilterOperators `json:"operator" validate:"required,valid_filter_operator"`

	// Value is the value to compare the column against.
	// This field is required.
	Value string `json:"value" validate:"required"`
}

// IncludeSelect defines the structure for specifying which columns to select
// from included related entities in a query.
type IncludeSelect struct {
	// The name of the relation to be included
	Name string `json:"name" validate:"required"`

	// The specific fields to select from the included relation/ columns
	Select []string `json:"select" validate:"required,dive,required"`
}

type CollectionResult[T any] struct {
	Total int64 `json:"total"`
	Items []*T  `json:"items"`
}
