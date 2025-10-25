package collectionquery

import (
	"fmt"
	"strings"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type QueryConstructor[T any] struct{}

func (qc *QueryConstructor[T]) ConstructQuery(
	db *gorm.DB,
	query CollectionQuery,
	withDelete bool,
) *gorm.DB {
	var t T
	qb := db.Model(&t)

	if withDelete {
		qb = qb.Unscoped()
	}

	// Parse schema for metadata
	// Use a cache as required by schema.Parse
	cache := &sync.Map{}
	sch, err := schema.Parse(&t, cache, db.NamingStrategy)
	if err != nil {
		return qb
	}

	tableName := sch.Table

	// Remove empty filters
	query = qc.removeEmptyFilter(query)

	// Apply Select
	if len(query.Select) > 0 {
		selectCols := make([]string, len(query.Select))
		for i, col := range query.Select {
			selectCols[i] = fmt.Sprintf("%s.%s", tableName, col)
		}
		qb = qb.Select(selectCols)
	}

	// Apply Where Clauses
	qb = qc.applyWhereConditions(tableName, qb, query.Where)

	// Apply GroupBy
	if len(query.GroupBy) > 0 {
		for _, col := range query.GroupBy {
			qb = qb.Group(fmt.Sprintf("%s.%s", tableName, col))
		}
	}

	// Apply HAVING
	if len(query.Having) > 0 {
		qb = qc.applyHavingConditions(tableName, qb, query.Having)
	}

	// Apply ORDER BY
	for _, order := range query.OrderBy {
		qb = qc.applyOrder(tableName, qb, order)
	}

	// Apply SKIP (offset)
	if query.Skip != nil {
		qb = qb.Offset(*query.Skip)
	}

	// Apply TAKE (limit)
	if query.Take != nil {
		qb = qb.Limit(*query.Take)
	}

	// Apply INCLUDES (preloading with joins)
	for _, include := range query.Includes {
		qb = qc.applyInclude(tableName, qb, include)
	}

	// Apply INCLUDE AND SELECT
	for _, includeSelect := range query.IncludeAndSelect {
		qb = qc.applyIncludeAndSelect(tableName, qb, includeSelect)
	}

	return qb
}

// hasField checks if schema has a field with the given database name
func (qc *QueryConstructor[T]) hasField(sch *schema.Schema, dbName string) bool {
	for _, field := range sch.Fields {
		if field.DBName == dbName {
			return true
		}
	}
	return false
}

// RemoveEmptyFilter removes empty where clause arrays
func (qc *QueryConstructor[T]) removeEmptyFilter(query CollectionQuery) CollectionQuery {
	filtered := make([][]Where, 0)
	for _, clause := range query.Where {
		if len(clause) > 0 {
			filtered = append(filtered, clause)
		}
	}

	query.Where = filtered

	havingFiltered := make([][]Where, 0)
	for _, clause := range query.Having {
		if len(clause) > 0 {
			havingFiltered = append(havingFiltered, clause)
		}
	}

	query.Having = havingFiltered
	return query
}

// RemoveFilter removes filters for a specific column
func (qc *QueryConstructor[T]) removeFilter(query CollectionQuery, key string) CollectionQuery {
	for i, whereGroup := range query.Where {
		filtered := make([]Where, 0)
		for _, clause := range whereGroup {
			if clause.Column != key {
				filtered = append(filtered, clause)
			}
		}
		query.Where[i] = filtered
	}

	for i, havingGroup := range query.Having {
		filtered := make([]Where, 0)
		for _, clause := range havingGroup {
			if clause.Column != key {
				filtered = append(filtered, clause)
			}
		}
		query.Having[i] = filtered
	}
	return query
}

// applyWhereConditions applies WHERE clauses with OR/AND logic
func (qc *QueryConstructor[T]) applyWhereConditions(
	tableName string,
	qb *gorm.DB,
	whereClauses [][]Where,
) *gorm.DB {
	for _, orGroup := range whereClauses {
		if len(orGroup) == 0 {
			continue
		}

		// Build OR conditions within this AND group
		qb = qb.Where(func(tx *gorm.DB) *gorm.DB {
			for i, clause := range orGroup {
				condition, args := qc.buildFilterCondition(tableName, clause)

				if i == 0 {
					tx = tx.Where(condition, args...)
				} else {
					tx = tx.Or(condition, args...)
				}
			}
			return tx
		})
	}
	return qb
}

// buildFilterCondition builds the WHERE condition string and returns args
func (qc *QueryConstructor[T]) buildFilterCondition(tableName string, clause Where) (string, []interface{}) {
	column := clause.Column
	op := clause.Operator
	value := clause.Value

	var queryCondition string

	if op == ArrayFilter {
		if strings.Contains(column, "->>") {
			parts := strings.SplitN(column, "->>", 2)
			mainColumn := parts[0]
			nestedColumn := parts[1]
			queryCondition = fmt.Sprintf(`("%s"."%s"->>'%s')::jsonb @> ?`, tableName, mainColumn, nestedColumn)
		} else {
			queryCondition = fmt.Sprintf(`"%s"."%s" @> ?`, tableName, column)
		}

		return queryCondition, []interface{}{value}
	}

	// Handler relation columns
	if strings.Contains(column, ".") {
		parts := strings.SplitN(column, ".", 2)
		relation := parts[0]
		field := parts[1]
		if strings.Contains(field, "->>") {
			fieldParts := strings.SplitN(field, "->>", 2)
			mainColumn := fieldParts[0]
			nestedColumn := fieldParts[1]
			queryCondition = fmt.Sprintf(`"%s"."%s" ->> '%s'`, relation, mainColumn, nestedColumn)
		} else if strings.Contains(field, "@>") {
			fieldParts := strings.SplitN(field, "@>", 2)
			mainColumn := fieldParts[0]
			queryCondition = fmt.Sprintf(`"%s"."%s" @>`, relation, mainColumn)
		} else {
			queryCondition = fmt.Sprintf(`"%s"."%s"`, relation, field)
		}
		return qc.applyOperators(queryCondition, op, value)
	}

	// Handle @> operator for main entity
	if strings.Contains(column, "@>") {
		parts := strings.SplitN(column, "@>", 2)
		mainColumn := parts[0]
		queryCondition = fmt.Sprintf(`"%s"."%s" @>`, tableName, mainColumn)
		return qc.applyOperators(queryCondition, op, value)
	}

	// Handle JSON field queries (->>)
	if strings.Contains(column, "->>") {
		parts := strings.SplitN(column, "->>", 2)
		mainColumn := parts[0]
		nestedColumn := parts[1]
		if strings.Contains(mainColumn, "->") {
			subParts := strings.SplitN(mainColumn, "->", 2)
			main := subParts[0]
			path := subParts[1]
			queryCondition = fmt.Sprintf(`"%s"."%s" -> '%s' ->> '%s'`, tableName, main, path, path)
		} else {
			queryCondition = fmt.Sprintf(`"%s"."%s" ->> '%s'`, tableName, mainColumn, nestedColumn)
		}
		return qc.applyOperators(queryCondition, op, value)
	}

	// Handle regular columns
	queryCondition = fmt.Sprintf(`"%s"."%s"`, tableName, column)
	return qc.applyOperators(queryCondition, op, value)
}

var operatorFormat = map[FilterOperators]string{
	EqualTo:              "%s = ?",
	GreaterThan:          "%s > ?",
	LessThan:             "%s < ?",
	GreaterThanOrEqualTo: "%s >= ?",
	LessThanOrEqualTo:    "%s <= ?",
	All:                  "%s = ALL(?)",
	Any:                  "%s = ANY(?)",
	Like:                 "%s LIKE ?",
	ILike:                "%s ILIKE ?",
	ArrayFilter:          "%s @> ?",
	ArrayContains:        "%s @> ?",
}

func (qc *QueryConstructor[T]) applyOperators(
	queryCondition string,
	op FilterOperators,
	value string,
) (string, []interface{}) {
	if format, ok := operatorFormat[op]; ok {
		var arg interface{} = value
		if op == Like || op == ILike {
			arg = fmt.Sprintf("%%%s%%", value)
		}
		return fmt.Sprintf(format, queryCondition), []interface{}{arg}
	}

	switch op {
	case Between:
		parts := strings.Split(value, ",")
		if len(parts) == 2 {
			return fmt.Sprintf("%s BETWEEN ? AND ?", queryCondition), []interface{}{parts[0], parts[1]}
		}
		return fmt.Sprintf("%s = ?", queryCondition), []interface{}{value}
	case In, NotIn:
		values := strings.Split(value, ",")
		operator := "IN"
		if op == NotIn {
			operator = "NOT IN"
		}
		return fmt.Sprintf("%s %s (?)", queryCondition, operator), []interface{}{values}
	case IsNull:
		return fmt.Sprintf("%s IS NULL", queryCondition), nil
	case IsNotNull:
		return fmt.Sprintf("%s IS NOT NULL", queryCondition), nil
	default:
		return fmt.Sprintf("%s %s ?", queryCondition, op), []interface{}{value}
	}
}

func (qc *QueryConstructor[T]) applyHavingConditions(
	tableName string,
	db *gorm.DB,
	conditions [][]Where,
) *gorm.DB {
	for _, orGroup := range conditions {
		if len(orGroup) == 0 {
			continue
		}

		db = db.Having(func(tx *gorm.DB) *gorm.DB {
			for i, clause := range orGroup {
				condition, args := qc.buildHavingCondition(tableName, clause)
				if i == 0 {
					tx = tx.Having(condition, args...)
				} else {
					tx = tx.Or(condition, args...)
				}
			}
			return tx
		})
	}
	return db
}

// buildHavingCondition builds HAVING condition with COUNT
func (qc *QueryConstructor[T]) buildHavingCondition(
	tableName string,
	clause Where,
) (string, []interface{}) {
	column := clause.Column
	op := clause.Operator
	value := clause.Value

	switch op {
	case Between:
		parts := strings.Split(value, ",")
		if len(parts) == 2 {
			return fmt.Sprintf("COUNT(%s.*) BETWEEN ? AND ?", tableName), []interface{}{parts[0], parts[1]}
		}
		return fmt.Sprintf("COUNT(%s.*) = ?", tableName), []interface{}{value}
	case In:
		values := strings.Split(value, ",")
		return fmt.Sprintf("COUNT(%s.*) IN ?", tableName), []interface{}{values}
	case Like:
		return fmt.Sprintf("COUNT(%s.*) LIKE ?", tableName), []interface{}{fmt.Sprintf("%%%s%%", value)}
	case Any, All:
		return fmt.Sprintf("COUNT(%s.%s) %s ?", tableName, column, op), []interface{}{value}
	default:
		return fmt.Sprintf("COUNT(%s.%s) %s ?", tableName, column, op), []interface{}{value}
	}
}

// applyOrder applies ordering
func (qc *QueryConstructor[T]) applyOrder(
	tableName string,
	db *gorm.DB,
	order Order,
) *gorm.DB {
	column := order.Column
	var orderStr string

	if strings.Contains(column, ".") {
		parts := strings.SplitN(column, ".", 2)
		orderStr = fmt.Sprintf(`"%s"."%s"`, parts[0], parts[1])
	} else {
		orderStr = fmt.Sprintf(`"%s"."%s"`, tableName, column)
	}

	if order.Direction != nil {
		orderStr += " " + string(*order.Direction)
	} else {
		orderStr += " ASC"
	}

	if order.Nulls != nil {
		switch *order.Nulls {
		case NullsFirst:
			orderStr += " NULLS FIRST"
		case NullsLast:
			orderStr += " NULLS LAST"
		}
	}

	return db.Order(orderStr)
}

func (qc *QueryConstructor[T]) applyInclude(
	tableName string,
	db *gorm.DB,
	include string,
) *gorm.DB {
	if strings.Contains(include, ".") {
		parts := strings.SplitN(include, ".", 2)
		parent := parts[0]
		child := parts[1]

		return db.Joins(parent).Joins(fmt.Sprintf("%s.%s", parent, child))
	}
	return db.Joins(include)
}

func (qc *QueryConstructor[T]) applyIncludeAndSelect(
	tableName string,
	db *gorm.DB,
	includeSelect IncludeSelect,
) *gorm.DB {
	relationName := includeSelect.Name

	// Use GORM's Joins with a callback to select specific columns
	// GORM will automatically handle the foreign key relationships
	return db.Joins(relationName, func(db *gorm.DB) *gorm.DB {
		// Build the select clause for the relation's columns
		selectCols := make([]string, len(includeSelect.Select))
		for i, field := range includeSelect.Select {
			selectCols[i] = fmt.Sprintf("%s.%s", relationName, field)
		}
		return db.Select(selectCols)
	})
}

// Find executes the query and returns results
func (qc *QueryConstructor[T]) Find(
	db *gorm.DB,
	query CollectionQuery,
	withDelete bool,
) (*CollectionResult[T], error) {
	// Build the base query once
	qb := qc.ConstructQuery(db, query, withDelete)

	// If only count is requested, return early
	if query.Count != nil && *query.Count {
		var count int64
		if err := qb.Count(&count).Error; err != nil {
			return nil, err
		}
		return &CollectionResult[T]{
			Total: count,
			Items: nil,
		}, nil
	}

	// Get total count without pagination (reuse same base query)
	var total int64
	countQuery := qb.Session(&gorm.Session{}) // Create a new session to avoid mutation
	countQuery = countQuery.Limit(-1).Offset(-1)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// Fetch the data with pagination applied
	var items []*T
	if err := qb.Find(&items).Error; err != nil {
		return nil, err
	}

	return &CollectionResult[T]{
		Total: total,
		Items: items,
	}, nil
}
