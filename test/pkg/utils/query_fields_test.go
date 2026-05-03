package utils_test

import (
	"testing"

	"go-api/pkg/utils"
)

func TestGetTableQueryConfig(t *testing.T) {
	cfg, ok := utils.GetTableQueryConfig(utils.TableProducts)
	if !ok {
		t.Fatalf("expected products config to exist")
	}
	if cfg.Table != "products" {
		t.Fatalf("unexpected table name: %s", cfg.Table)
	}
	if len(cfg.Fields) == 0 {
		t.Fatalf("expected non-empty fields map")
	}
	if len(cfg.DefaultSorts) == 0 {
		t.Fatalf("expected at least one default sort")
	}
}

func TestGetTableQueryConfigMissing(t *testing.T) {
	_, ok := utils.GetTableQueryConfig("nonexistent_table")
	if ok {
		t.Fatalf("expected missing config for unknown table")
	}
}

func TestKnownTableConstants(t *testing.T) {
	tables := []string{
		utils.TableProducts,
		utils.TableModules,
		utils.TableActions,
		utils.TableResources,
		utils.TableRoles,
		utils.TableUsers,
	}
	for _, tbl := range tables {
		if _, ok := utils.GetTableQueryConfig(tbl); !ok {
			t.Fatalf("missing config for table constant %q", tbl)
		}
	}
}
