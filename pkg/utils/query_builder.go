package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	basev1 "github.com/lamquangmanh/protobuf/gen/go/proto/base/v1"
)

// NormalizeSorts validates sort fields against a table config and returns normalized SQL sort clauses.
// Sample input:  sorts=[{Field:"name",Order:ASC}], cfg=products, maxSortFields=5
// Sample output: [{Column:"name", Order:"ASC"}], nil
func NormalizeSorts(sorts []*basev1.Sort, cfg QueryTableConfig, maxSortFields int) ([]SortClause, error) {
	if maxSortFields > 0 && len(sorts) > maxSortFields {
		return nil, fmt.Errorf("sort fields exceed max allowed: %d", maxSortFields)
	}

	result := make([]SortClause, 0, len(sorts))
	usedColumns := map[string]bool{}
	for _, sort := range sorts {
		if sort == nil {
			continue
		}
		field := strings.TrimSpace(sort.GetField())
		fieldCfg, ok := cfg.Fields[field]
		if !ok {
			return nil, fmt.Errorf("invalid sort field: %s", field)
		}
		if !fieldCfg.Sortable {
			return nil, fmt.Errorf("field is not sortable: %s", field)
		}
		if usedColumns[fieldCfg.Column] {
			continue
		}

		order := "DESC"
		switch sort.GetOrder() {
		case basev1.SortOrder_SORT_ORDER_UNSPECIFIED, basev1.SortOrder_SORT_ORDER_DESC:
			order = "DESC"
		case basev1.SortOrder_SORT_ORDER_ASC:
			order = "ASC"
		default:
			return nil, fmt.Errorf("invalid sort order for field: %s", field)
		}

		result = append(result, SortClause{Column: fieldCfg.Column, Order: order})
		usedColumns[fieldCfg.Column] = true
	}

	if len(result) == 0 {
		result = append(result, cfg.DefaultSorts...)
	}
	return result, nil
}

// NormalizeFilters validates filters against a table config and converts proto values into typed SQL args.
// Sample input:  filters=[{Field:"productId",Operator:IN,StringValues:["uuid-1","uuid-2"]}]
// Sample output: [{Column:"product_id",DataType:UUID,Operator:IN,Value:[]uuid.UUID{...}}], nil
func NormalizeFilters(filters []*basev1.Filter, cfg QueryTableConfig, maxFilterFields int) ([]FilterClause, error) {
	if maxFilterFields > 0 && len(filters) > maxFilterFields {
		return nil, fmt.Errorf("filter fields exceed max allowed: %d", maxFilterFields)
	}

	result := make([]FilterClause, 0, len(filters))
	for _, filter := range filters {
		if filter == nil {
			continue
		}

		field := strings.TrimSpace(filter.GetField())
		fieldCfg, ok := cfg.Fields[field]
		if !ok {
			return nil, fmt.Errorf("invalid filter field: %s", field)
		}

		operator := filter.GetOperator()
		if operator == basev1.FilterOperator_FILTER_OPERATOR_UNSPECIFIED {
			operator = basev1.FilterOperator_FILTER_OPERATOR_EQUAL
		}
		if _, ok := fieldCfg.AllowedOperators[operator]; !ok {
			return nil, fmt.Errorf("operator is not allowed for field %s", field)
		}

		value, err := ParseFilterValue(filter, fieldCfg.DataType, operator)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for field %s: %w", field, err)
		}

		result = append(result, FilterClause{
			Column:   fieldCfg.Column,
			DataType: fieldCfg.DataType,
			Operator: operator,
			Value:    value,
		})
	}

	return result, nil
}

// BuildListAndCountQueries builds SQL statements and argument lists for list and count queries.
// Sample input:  cfg.Table="products", sorts=[name ASC], filters=[name = "demo"], limit=10, offset=0
// Sample output: listSQL="SELECT ... FROM products ... ORDER BY name ASC LIMIT $2 OFFSET $3",
//
//	listArgs=["demo",10,0], countSQL="SELECT COUNT(*) FROM products ...", countArgs=["demo"]
func BuildListAndCountQueries(cfg QueryTableConfig, sorts []SortClause, filters []FilterClause, limit int32, offset int32) (string, []any, string, []any) {
	whereParts := make([]string, 0, len(filters)+1)
	whereParts = append(whereParts, fmt.Sprintf("%s IS NULL", cfg.SoftDeleteColumn))

	args := make([]any, 0, len(filters)+2)
	argIndex := 1
	for _, filter := range filters {
		condition := buildFilterCondition(filter, argIndex)
		whereParts = append(whereParts, condition)
		args = append(args, filter.Value)
		argIndex++
	}

	whereClause := strings.Join(whereParts, " AND ")
	orderParts := make([]string, 0, len(sorts))
	for _, sort := range sorts {
		orderParts = append(orderParts, fmt.Sprintf("%s %s", sort.Column, sort.Order))
	}
	if len(orderParts) == 0 {
		orderParts = append(orderParts, "created_at DESC")
	}
	orderClause := strings.Join(orderParts, ", ")

	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", cfg.Table, whereClause)
	countArgs := append([]any(nil), args...)

	listSQL := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT $%d OFFSET $%d",
		strings.Join(cfg.SelectColumns, ", "),
		cfg.Table,
		whereClause,
		orderClause,
		argIndex,
		argIndex+1,
	)
	listArgs := append(args, limit, offset)

	return listSQL, listArgs, countSQL, countArgs
}

// buildFilterCondition builds one SQL predicate for a normalized filter clause.
// Sample input:  filter={Column:"name",Operator:LIKE}, argPos=1
// Sample output: "name ILIKE '%' || $1 || '%'"
func buildFilterCondition(filter FilterClause, argPos int) string {
	placeholder := fmt.Sprintf("$%d", argPos)
	switch filter.Operator {
	case basev1.FilterOperator_FILTER_OPERATOR_EQUAL:
		return fmt.Sprintf("%s = %s", filter.Column, placeholder)
	case basev1.FilterOperator_FILTER_OPERATOR_NOT_EQUAL:
		return fmt.Sprintf("%s <> %s", filter.Column, placeholder)
	case basev1.FilterOperator_FILTER_OPERATOR_GREATER_THAN:
		return fmt.Sprintf("%s > %s", filter.Column, placeholder)
	case basev1.FilterOperator_FILTER_OPERATOR_LESS_THAN:
		return fmt.Sprintf("%s < %s", filter.Column, placeholder)
	case basev1.FilterOperator_FILTER_OPERATOR_GREATER_THAN_OR_EQUAL:
		return fmt.Sprintf("%s >= %s", filter.Column, placeholder)
	case basev1.FilterOperator_FILTER_OPERATOR_LESS_THAN_OR_EQUAL:
		return fmt.Sprintf("%s <= %s", filter.Column, placeholder)
	case basev1.FilterOperator_FILTER_OPERATOR_LIKE:
		return fmt.Sprintf("%s ILIKE '%%' || %s || '%%'", filter.Column, placeholder)
	case basev1.FilterOperator_FILTER_OPERATOR_IN:
		return fmt.Sprintf("%s = ANY(%s)", filter.Column, placeholder)
	case basev1.FilterOperator_FILTER_OPERATOR_NOT_IN:
		return fmt.Sprintf("NOT (%s = ANY(%s))", filter.Column, placeholder)
	default:
		return fmt.Sprintf("%s = %s", filter.Column, placeholder)
	}
}

// ParseFilterValue parses filter values from proto into the correct Go type for SQL binding.
// Sample input:  dataType=FieldDataTypeUUID, operator=IN, StringValues=["550e8400-e29b-41d4-a716-446655440000"]
// Sample output: []uuid.UUID{550e8400-e29b-41d4-a716-446655440000}, nil
func ParseFilterValue(filter *basev1.Filter, dataType FieldDataType, operator basev1.FilterOperator) (any, error) {
	isArrayOperator := operator == basev1.FilterOperator_FILTER_OPERATOR_IN || operator == basev1.FilterOperator_FILTER_OPERATOR_NOT_IN

	singleString := strings.TrimSpace(filter.GetStringValue())
	stringValues := filter.GetStringValues()

	switch dataType {
	case FieldDataTypeText:
		if isArrayOperator {
			if len(stringValues) == 0 && singleString != "" {
				return []string{singleString}, nil
			}
			result := make([]string, 0, len(stringValues))
			for _, value := range stringValues {
				v := strings.TrimSpace(value)
				if v == "" {
					continue
				}
				result = append(result, v)
			}
			if len(result) == 0 {
				return nil, fmt.Errorf("empty values")
			}
			return result, nil
		}
		if singleString == "" && len(stringValues) > 0 {
			singleString = strings.TrimSpace(stringValues[0])
		}
		if singleString == "" {
			return nil, fmt.Errorf("empty value")
		}
		return singleString, nil

	case FieldDataTypeUUID:
		if isArrayOperator {
			values := stringValues
			if len(values) == 0 && singleString != "" {
				values = []string{singleString}
			}
			result := make([]uuid.UUID, 0, len(values))
			for _, value := range values {
				v := strings.TrimSpace(value)
				if v == "" {
					continue
				}
				parsed, err := uuid.Parse(v)
				if err != nil {
					return nil, err
				}
				result = append(result, parsed)
			}
			if len(result) == 0 {
				return nil, fmt.Errorf("empty values")
			}
			return result, nil
		}
		if singleString == "" && len(stringValues) > 0 {
			singleString = strings.TrimSpace(stringValues[0])
		}
		if singleString == "" {
			return nil, fmt.Errorf("empty value")
		}
		return uuid.Parse(singleString)

	case FieldDataTypeTimestamp:
		if singleString == "" && len(stringValues) > 0 {
			singleString = strings.TrimSpace(stringValues[0])
		}
		if singleString == "" {
			return nil, fmt.Errorf("empty value")
		}
		parsed, err := time.Parse(time.RFC3339, singleString)
		if err != nil {
			return nil, err
		}
		return parsed, nil

	default:
		return nil, fmt.Errorf("unsupported data type")
	}
}
