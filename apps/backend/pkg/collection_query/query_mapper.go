package collectionquery

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// queryEncoder is a helper struct to manage encoding process.
type queryEncoder struct {
	query       *CollectionQuery
	queryParams []string
}

// Function to encode a CollectionQuery object to a custom URL query string
func EncodeColllectionQuery(query CollectionQuery) string {
	encoder := &queryEncoder{
		query:       &query,
		queryParams: make([]string, 0, 9),
	}

	encoder.encodeSelect()
	encoder.encodeWhere()
	encoder.encodeTake()
	encoder.encodeSkip()
	encoder.encodeOrderBy()
	encoder.encodeIncludes()
	encoder.encodeGroupBy()
	encoder.encodeHaving()
	encoder.encodeCount()

	return strings.Join(encoder.queryParams, "&")
}

func (e *queryEncoder) encodeSelect() {
	if len(e.query.Select) > 0 {
		e.queryParams = append(e.queryParams, fmt.Sprintf("s=%s", strings.Join(e.query.Select, ",")))
	}
}

func (e *queryEncoder) encodeWhere() {
	if len(e.query.Where) > 0 {
		e.queryParams = append(e.queryParams, fmt.Sprintf("w=%s", encodeWhere(e.query.Where)))
	}
}

func (e *queryEncoder) encodeTake() {
	if e.query.Take != nil {
		e.queryParams = append(e.queryParams, fmt.Sprintf("t=%d", *e.query.Take))
	}
}

func (e *queryEncoder) encodeSkip() {
	if e.query.Skip != nil {
		e.queryParams = append(e.queryParams, fmt.Sprintf("sk=%d", *e.query.Skip))
	}
}

func (e *queryEncoder) encodeOrderBy() {
	if len(e.query.OrderBy) > 0 {
		e.queryParams = append(e.queryParams, fmt.Sprintf("o=%s", encodeOrderBy(e.query.OrderBy)))
	}
}

func (e *queryEncoder) encodeIncludes() {
	if len(e.query.Includes) > 0 {
		e.queryParams = append(e.queryParams, fmt.Sprintf("i=%s", strings.Join(e.query.Includes, ",")))
	}
}

func (e *queryEncoder) encodeGroupBy() {
	if len(e.query.GroupBy) > 0 {
		e.queryParams = append(e.queryParams, fmt.Sprintf("g=%s", strings.Join(e.query.GroupBy, ",")))
	}
}

func (e *queryEncoder) encodeHaving() {
	if len(e.query.Having) > 0 {
		e.queryParams = append(e.queryParams, fmt.Sprintf("h=%s", encodeWhere(e.query.Having)))
	}
}

func (e *queryEncoder) encodeCount() {
	if e.query.Count != nil {
		e.queryParams = append(e.queryParams, fmt.Sprintf("c=%t", *e.query.Count))
	}
}

// queryDecoder is a helper struct to manage decoding process.
type queryDecoder struct {
	query       *CollectionQuery
	queryParams url.Values
}

func DecodeCollectionQuery(queryString string) (CollectionQuery, error) {
	if queryString == "" {
		return CollectionQuery{}, nil
	}

	queryParams, err := url.ParseQuery(queryString)
	if err != nil {
		return CollectionQuery{}, err
	}

	decoder := &queryDecoder{
		query:       &CollectionQuery{},
		queryParams: queryParams,
	}

	decoder.decodeSelect()
	decoder.decodeWhere()
	decoder.decodeTake()
	decoder.decodeSkip()
	decoder.decodeOrderBy()
	decoder.decodeIncludes()
	decoder.decodeGroupBy()
	decoder.decodeHaving()
	decoder.decodeCount()

	return *decoder.query, nil
}

func (d *queryDecoder) get(key string) (string, bool) {
	if d.queryParams.Has(key) && d.queryParams.Get(key) != "" {
		return d.queryParams.Get(key), true
	}
	return "", false
}

func (d *queryDecoder) decodeSelect() {
	if val, ok := d.get("s"); ok {
		d.query.Select = strings.Split(val, ",")
	}
}

func (d *queryDecoder) decodeWhere() {
	if val, ok := d.get("w"); ok {
		d.query.Where = decodeWhere(val)
	}
}

func (d *queryDecoder) decodeTake() {
	if val, ok := d.get("t"); ok {
		parsedInt, err := strconv.ParseInt(val, 10, 64)
		takeValue := 10
		if err == nil {
			takeValue = int(parsedInt)
			d.query.Take = &takeValue
		}
		d.query.Take = &takeValue
	}
}

func (d *queryDecoder) decodeSkip() {
	if val, ok := d.get("sk"); ok {
		parsedInt, err := strconv.ParseInt(val, 10, 64)
		skipValue := 0
		if err == nil {
			skipValue = int(parsedInt)
		}
		d.query.Skip = &skipValue
	}
}

func (d *queryDecoder) decodeOrderBy() {
	if val, ok := d.get("o"); ok {
		d.query.OrderBy = decodeOrderBy(val)
	}
}

func (d *queryDecoder) decodeIncludes() {
	if val, ok := d.get("i"); ok {
		d.query.Includes = strings.Split(val, ",")
	}
}

func (d *queryDecoder) decodeGroupBy() {
	if val, ok := d.get("g"); ok {
		d.query.GroupBy = strings.Split(val, ",")
	}
}

func (d *queryDecoder) decodeHaving() {
	if val, ok := d.get("h"); ok {
		d.query.Having = decodeWhere(val)
	}
}

func (d *queryDecoder) decodeCount() {
	if val, ok := d.get("c"); ok {
		if countValue, err := strconv.ParseBool(val); err == nil {
			d.query.Count = &countValue
		}
	}
}

func encodeWhere(where [][]Where) string {
	groups := make([]string, 0, len(where))
	for _, group := range where {
		groups = append(groups, encodeWhereGroup(group))
	}
	return strings.Join(groups, string(WhereAND))
}

func decodeWhere(encoded string) [][]Where {
	if encoded == "" {
		return [][]Where{}
	}

	encodedGroups := strings.Split(encoded, string(WhereAND))
	decodedWhereGroups := make([][]Where, 0, len(encodedGroups))

	for _, groupStr := range encodedGroups {
		if groupStr == "" {
			continue
		}
		decodedGroup := decodeWhereGroup(groupStr)
		decodedWhereGroups = append(decodedWhereGroups, decodedGroup)
	}

	return decodedWhereGroups
}

func encodeWhereGroup(group []Where) string {
	items := make([]string, 0, len(group))
	for _, item := range group {
		items = append(items, encodeWhereItem(item))
	}
	return strings.Join(items, string(WhereOR))
}

func decodeWhereGroup(encoded string) []Where {
	if encoded == "" {
		return []Where{}
	}

	encodedItems := strings.Split(encoded, string(WhereOR))
	decodeWhereGroup := make([]Where, len(encodedItems))

	for _, itemStr := range encodedItems {
		if itemStr == "" {
			continue
		}
		decodedItem := decodeWhereItem(itemStr)
		decodeWhereGroup = append(decodeWhereGroup, decodedItem)
	}
	return decodeWhereGroup
}

func encodeWhereItem(item Where) string {
	return fmt.Sprintf("%s%s%s%s%s", item.Column, string(WhereEqual), item.Operator, string(WhereEqual), item.Value)
}

func decodeWhereItem(encoded string) Where {
	parts := strings.Split(encoded, string(WhereEqual))
	return Where{
		Column:   parts[0],
		Operator: FilterOperators(parts[1]),
		Value:    parts[2],
	}
}

func encodeOrderBy(orderBy []Order) string {
	orders := make([]string, 0, len(orderBy))
	for _, order := range orderBy {
		orders = append(orders, encodeOrderItem(order))
	}
	return strings.Join(orders, string(OrderBy))
}

func decodeOrderBy(encoded string) []Order {
	if encoded == "" {
		return []Order{}
	}

	encodedItems := strings.Split(encoded, string(OrderBy))
	decodedOrders := make([]Order, 0, len(encodedItems))
	for _, itemStr := range encodedItems {
		if itemStr == "" {
			continue
		}
		decodedItem := decodeOrderItem(itemStr)
		decodedOrders = append(decodedOrders, decodedItem)
	}
	return decodedOrders
}

func encodeOrderItem(item Order) string {
	return fmt.Sprintf("%s%s%s",
		item.Column,
		OrderItem,
		valueOrDefaultString((*string)(item.Direction), "ASC"),
	)
}

func decodeOrderItem(encoded string) Order {
	parts := strings.Split(encoded, string(OrderItem))
	dir := SortDirection(parts[1])
	return Order{
		Column:    parts[0],
		Direction: &dir,
	}
}

// helper for string value or default value if value is nil
func valueOrDefaultString(value *string, defaultValue string) string {
	if value == nil {
		return defaultValue
	}
	return *value
}
