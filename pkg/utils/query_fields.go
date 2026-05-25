package utils

import basev1 "github.com/lamquangmanh/protobuf/gen/go/proto/base/v1"

// FieldDataType declares supported data kinds for filter value parsing.
type FieldDataType string

const (
	FieldDataTypeText      FieldDataType = "text"
	FieldDataTypeUUID      FieldDataType = "uuid"
	FieldDataTypeTimestamp FieldDataType = "timestamp"
)

type QueryField struct {
	Column           string
	DataType         FieldDataType
	Sortable         bool
	AllowedOperators map[basev1.FilterOperator]struct{}
}

// SortClause represents one validated ORDER BY segment.
type SortClause struct {
	Column string
	Order  string
}

// FilterClause represents one validated WHERE clause condition.
type FilterClause struct {
	Column   string
	DataType FieldDataType
	Operator basev1.FilterOperator
	Value    any
}

// QueryTableConfig describes sortable/filterable fields and default query behavior per table.
type QueryTableConfig struct {
	Table            string
	SelectColumns    []string
	SoftDeleteColumn string
	DefaultSorts     []SortClause
	Fields           map[string]QueryField
}

const (
	TableProducts  = "products"
	TableModules   = "modules"
	TableResources = "resources"
	TableActions   = "actions"
	TableRoles     = "roles"
	TableUsers     = "users"
)

// ops builds a set-like map from filter operators for quick membership checks.
// Sample input:  ops(EQUAL, LIKE)
// Sample output: map[EQUAL:{} LIKE:{}]
func ops(operators ...basev1.FilterOperator) map[basev1.FilterOperator]struct{} {
	result := make(map[basev1.FilterOperator]struct{}, len(operators))
	for _, operator := range operators {
		result[operator] = struct{}{}
	}
	return result
}

var tableQueryConfigs = map[string]QueryTableConfig{
	TableProducts: {
		Table:            "products",
		SelectColumns:    []string{"product_id", "name", "description", "url", "created_at", "created_user_id", "updated_at", "updated_user_id", "deleted_at", "deleted_user_id"},
		SoftDeleteColumn: "deleted_at",
		DefaultSorts:     []SortClause{{Column: "created_at", Order: "DESC"}},
		Fields: map[string]QueryField{
			"product_id":  {Column: "product_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_NOT_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN, basev1.FilterOperator_FILTER_OPERATOR_NOT_IN)},
			"productId":   {Column: "product_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_NOT_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN, basev1.FilterOperator_FILTER_OPERATOR_NOT_IN)},
			"name":        {Column: "name", DataType: FieldDataTypeText, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_NOT_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LIKE, basev1.FilterOperator_FILTER_OPERATOR_IN, basev1.FilterOperator_FILTER_OPERATOR_NOT_IN)},
			"description": {Column: "description", DataType: FieldDataTypeText, Sortable: false, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_NOT_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LIKE)},
			"url":         {Column: "url", DataType: FieldDataTypeText, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_NOT_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LIKE, basev1.FilterOperator_FILTER_OPERATOR_IN, basev1.FilterOperator_FILTER_OPERATOR_NOT_IN)},
			"created_at":  {Column: "created_at", DataType: FieldDataTypeTimestamp, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_NOT_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_GREATER_THAN, basev1.FilterOperator_FILTER_OPERATOR_GREATER_THAN_OR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LESS_THAN, basev1.FilterOperator_FILTER_OPERATOR_LESS_THAN_OR_EQUAL)},
			"createdAt":   {Column: "created_at", DataType: FieldDataTypeTimestamp, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_NOT_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_GREATER_THAN, basev1.FilterOperator_FILTER_OPERATOR_GREATER_THAN_OR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LESS_THAN, basev1.FilterOperator_FILTER_OPERATOR_LESS_THAN_OR_EQUAL)},
			"updated_at":  {Column: "updated_at", DataType: FieldDataTypeTimestamp, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_NOT_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_GREATER_THAN, basev1.FilterOperator_FILTER_OPERATOR_GREATER_THAN_OR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LESS_THAN, basev1.FilterOperator_FILTER_OPERATOR_LESS_THAN_OR_EQUAL)},
			"updatedAt":   {Column: "updated_at", DataType: FieldDataTypeTimestamp, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_NOT_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_GREATER_THAN, basev1.FilterOperator_FILTER_OPERATOR_GREATER_THAN_OR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LESS_THAN, basev1.FilterOperator_FILTER_OPERATOR_LESS_THAN_OR_EQUAL)},
		},
	},
	TableModules: {
		Table:            "modules",
		SelectColumns:    []string{"module_id", "name", "description", "url", "icon", "product_id", "created_at", "created_user_id", "updated_at", "updated_user_id", "deleted_at", "deleted_user_id"},
		SoftDeleteColumn: "deleted_at",
		DefaultSorts:     []SortClause{{Column: "created_at", Order: "DESC"}},
		Fields: map[string]QueryField{
			"module_id":  {Column: "module_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"moduleId":   {Column: "module_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"name":       {Column: "name", DataType: FieldDataTypeText, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LIKE)},
			"product_id": {Column: "product_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"productId":  {Column: "product_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
		},
	},
	TableResources: {
		Table:            "resources",
		SelectColumns:    []string{"resource_id", "name", "module_id", "created_at", "created_user_id", "updated_at", "updated_user_id", "deleted_at", "deleted_user_id"},
		SoftDeleteColumn: "deleted_at",
		DefaultSorts:     []SortClause{{Column: "created_at", Order: "DESC"}},
		Fields: map[string]QueryField{
			"resource_id": {Column: "resource_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"resourceId":  {Column: "resource_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"module_id":   {Column: "module_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"moduleId":    {Column: "module_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"name":        {Column: "name", DataType: FieldDataTypeText, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LIKE)},
		},
	},
	TableActions: {
		Table:            "actions",
		SelectColumns:    []string{"action_id", "resource_id", "name", "description", "request_type", "url", "method", "created_at", "created_user_id", "updated_at", "updated_user_id", "deleted_at", "deleted_user_id"},
		SoftDeleteColumn: "deleted_at",
		DefaultSorts:     []SortClause{{Column: "created_at", Order: "DESC"}},
		Fields: map[string]QueryField{
			"action_id":    {Column: "action_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"actionId":     {Column: "action_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"resource_id":  {Column: "resource_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"resourceId":   {Column: "resource_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"name":         {Column: "name", DataType: FieldDataTypeText, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LIKE)},
			"request_type": {Column: "request_type", DataType: FieldDataTypeText, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"requestType":  {Column: "request_type", DataType: FieldDataTypeText, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
		},
	},
	TableRoles: {
		Table:            "roles",
		SelectColumns:    []string{"role_id", "name", "description", "module_id", "created_at", "created_user_id", "updated_at", "updated_user_id", "deleted_at", "deleted_user_id"},
		SoftDeleteColumn: "deleted_at",
		DefaultSorts:     []SortClause{{Column: "created_at", Order: "DESC"}},
		Fields: map[string]QueryField{
			"role_id":     {Column: "role_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"roleId":      {Column: "role_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"module_id":   {Column: "module_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"moduleId":    {Column: "module_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"name":        {Column: "name", DataType: FieldDataTypeText, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LIKE)},
			"description": {Column: "description", DataType: FieldDataTypeText, Sortable: false, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LIKE)},
		},
	},
	TableUsers: {
		Table:            "users",
		SelectColumns:    []string{"user_id", "user_name", "email", "password", "phone", "avatar", "status", "created_at", "created_user_id", "updated_at", "updated_user_id", "deleted_at", "deleted_user_id"},
		SoftDeleteColumn: "deleted_at",
		DefaultSorts:     []SortClause{{Column: "created_at", Order: "DESC"}},
		Fields: map[string]QueryField{
			"user_id":   {Column: "user_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"userId":    {Column: "user_id", DataType: FieldDataTypeUUID, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
			"user_name": {Column: "user_name", DataType: FieldDataTypeText, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LIKE)},
			"userName":  {Column: "user_name", DataType: FieldDataTypeText, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LIKE)},
			"email":     {Column: "email", DataType: FieldDataTypeText, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_LIKE)},
			"status":    {Column: "status", DataType: FieldDataTypeText, Sortable: true, AllowedOperators: ops(basev1.FilterOperator_FILTER_OPERATOR_EQUAL, basev1.FilterOperator_FILTER_OPERATOR_IN)},
		},
	},
}

// GetTableQueryConfig returns query configuration for a logical table key.
// Sample input:  table = TableProducts
// Sample output: cfg.Table == "products", ok == true
func GetTableQueryConfig(table string) (QueryTableConfig, bool) {
	cfg, ok := tableQueryConfigs[table]
	return cfg, ok
}
