package utils_test

import (
	"strings"
	"testing"
	"time"

	basepb "go-api/pkg/api/basepb"
	"go-api/pkg/utils"

	"github.com/google/uuid"
)

func strPtr(v string) *string { return &v }

func mustTableConfig(t *testing.T) utils.QueryTableConfig {
	t.Helper()
	cfg, ok := utils.GetTableQueryConfig(utils.TableProducts)
	if !ok {
		t.Fatalf("products table config not found")
	}
	return cfg
}

func TestNormalizeSorts(t *testing.T) {
	cfg := mustTableConfig(t)

	t.Run("default sort when empty", func(t *testing.T) {
		got, err := utils.NormalizeSorts(nil, cfg, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 || got[0].Column != "created_at" || got[0].Order != "DESC" {
			t.Fatalf("unexpected default sorts: %+v", got)
		}
	})

	t.Run("deduplicate by normalized column", func(t *testing.T) {
		sorts := []*basepb.Sort{
			{Field: "name", Order: basepb.SortOrder_SORT_ORDER_ASC},
			{Field: "name", Order: basepb.SortOrder_SORT_ORDER_DESC},
			{Field: "updatedAt", Order: basepb.SortOrder_SORT_ORDER_DESC},
		}
		got, err := utils.NormalizeSorts(sorts, cfg, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("expected 2 sorts after dedup, got %d", len(got))
		}
		if got[0].Column != "name" || got[0].Order != "ASC" {
			t.Fatalf("unexpected first sort: %+v", got[0])
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		_, err := utils.NormalizeSorts([]*basepb.Sort{{Field: "bad_field"}}, cfg, 5)
		if err == nil || !strings.Contains(err.Error(), "invalid sort field") {
			t.Fatalf("expected invalid sort field error, got: %v", err)
		}
	})

	t.Run("max sort limit", func(t *testing.T) {
		sorts := []*basepb.Sort{{Field: "name"}, {Field: "updated_at"}}
		_, err := utils.NormalizeSorts(sorts, cfg, 1)
		if err == nil || !strings.Contains(err.Error(), "sort fields exceed max allowed") {
			t.Fatalf("expected max sort error, got: %v", err)
		}
	})
}

func TestNormalizeFilters(t *testing.T) {
	cfg := mustTableConfig(t)
	id1 := uuid.New().String()
	id2 := uuid.New().String()

	t.Run("text and uuid filters", func(t *testing.T) {
		filters := []*basepb.Filter{
			{Field: "name", Operator: basepb.FilterOperator_FILTER_OPERATOR_LIKE, StringValue: strPtr("abc")},
			{Field: "productId", Operator: basepb.FilterOperator_FILTER_OPERATOR_IN, StringValues: []string{id1, id2}},
		}
		got, err := utils.NormalizeFilters(filters, cfg, 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("expected 2 filters, got %d", len(got))
		}
		if got[0].Column != "name" || got[0].Value != "abc" {
			t.Fatalf("unexpected first filter: %+v", got[0])
		}
		uuids, ok := got[1].Value.([]uuid.UUID)
		if !ok || len(uuids) != 2 {
			t.Fatalf("expected []uuid.UUID with len=2, got %T %+v", got[1].Value, got[1].Value)
		}
	})

	t.Run("invalid operator for field", func(t *testing.T) {
		filters := []*basepb.Filter{{Field: "description", Operator: basepb.FilterOperator_FILTER_OPERATOR_IN, StringValues: []string{"x"}}}
		_, err := utils.NormalizeFilters(filters, cfg, 10)
		if err == nil || !strings.Contains(err.Error(), "operator is not allowed") {
			t.Fatalf("expected operator not allowed error, got: %v", err)
		}
	})

	t.Run("invalid uuid value", func(t *testing.T) {
		filters := []*basepb.Filter{{Field: "product_id", Operator: basepb.FilterOperator_FILTER_OPERATOR_EQUAL, StringValue: strPtr("not-uuid")}}
		_, err := utils.NormalizeFilters(filters, cfg, 10)
		if err == nil || !strings.Contains(err.Error(), "invalid filter value") {
			t.Fatalf("expected invalid filter value error, got: %v", err)
		}
	})

	t.Run("max filter limit", func(t *testing.T) {
		filters := []*basepb.Filter{{Field: "name", StringValue: strPtr("a")}, {Field: "name", StringValue: strPtr("b")}}
		_, err := utils.NormalizeFilters(filters, cfg, 1)
		if err == nil || !strings.Contains(err.Error(), "filter fields exceed max allowed") {
			t.Fatalf("expected max filter error, got: %v", err)
		}
	})
}

func TestBuildListAndCountQueries(t *testing.T) {
	cfg := mustTableConfig(t)
	sorts := []utils.SortClause{{Column: "name", Order: "ASC"}}
	filters := []utils.FilterClause{{Column: "name", DataType: utils.FieldDataTypeText, Operator: basepb.FilterOperator_FILTER_OPERATOR_EQUAL, Value: "demo"}}

	listSQL, listArgs, countSQL, countArgs := utils.BuildListAndCountQueries(cfg, sorts, filters, 10, 20)
	if !strings.Contains(listSQL, "FROM products") || !strings.Contains(listSQL, "ORDER BY name ASC") {
		t.Fatalf("unexpected list sql: %s", listSQL)
	}
	if !strings.Contains(countSQL, "SELECT COUNT(*) FROM products") {
		t.Fatalf("unexpected count sql: %s", countSQL)
	}
	if len(listArgs) != 3 {
		t.Fatalf("expected 3 list args, got %d", len(listArgs))
	}
	if len(countArgs) != 1 {
		t.Fatalf("expected 1 count arg, got %d", len(countArgs))
	}
	listArgs[0] = "changed"
	if countArgs[0] != "demo" {
		t.Fatalf("count args should be independent copy")
	}
}

func TestParseFilterValueTimestamp(t *testing.T) {
	v := "2026-04-25T10:00:00Z"
	out, err := utils.ParseFilterValue(
		&basepb.Filter{Field: "created_at", StringValue: &v},
		utils.FieldDataTypeTimestamp,
		basepb.FilterOperator_FILTER_OPERATOR_EQUAL,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out.(time.Time); !ok {
		t.Fatalf("expected time.Time, got %T", out)
	}
}
