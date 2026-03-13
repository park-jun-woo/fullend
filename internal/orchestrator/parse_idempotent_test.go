package orchestrator

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/geul-org/fullend/internal/funcspec"
	"github.com/geul-org/fullend/internal/policy"
	"github.com/geul-org/fullend/internal/projectconfig"
	"github.com/geul-org/fullend/internal/statemachine"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
	stmlparser "github.com/geul-org/fullend/internal/stml/parser"
)

// findSpecsDir locates specs/gigbridge relative to the project root.
func findSpecsDir(t *testing.T) string {
	t.Helper()
	// Walk up from this test file to find project root (where go.mod lives).
	_, thisFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("cannot find project root (go.mod)")
		}
		dir = parent
	}
	specsDir := filepath.Join(dir, "specs", "gigbridge")
	if _, err := os.Stat(specsDir); err != nil {
		t.Skipf("specs/gigbridge not found: %v", err)
	}
	return specsDir
}

// TestParseIdempotency verifies that parsing each SSOT twice with the same
// input produces identical results. This is a prerequisite for the parse
// consolidation in Phase017.
func TestParseIdempotency(t *testing.T) {
	specsDir := findSpecsDir(t)

	t.Run("OpenAPI", func(t *testing.T) {
		path := filepath.Join(specsDir, "api", "openapi.yaml")
		doc1, err1 := openapi3.NewLoader().LoadFromFile(path)
		doc2, err2 := openapi3.NewLoader().LoadFromFile(path)
		if err1 != nil || err2 != nil {
			t.Fatalf("load errors: %v / %v", err1, err2)
		}
		if !reflect.DeepEqual(doc1, doc2) {
			t.Error("OpenAPI: two parses produced different results")
		}
	})

	t.Run("SymbolTable", func(t *testing.T) {
		st1, err1 := ssacvalidator.LoadSymbolTable(specsDir)
		st2, err2 := ssacvalidator.LoadSymbolTable(specsDir)
		if err1 != nil || err2 != nil {
			t.Fatalf("load errors: %v / %v", err1, err2)
		}
		if !reflect.DeepEqual(st1, st2) {
			diffSymbolTables(t, st1, st2)
		}
	})

	t.Run("SSaC", func(t *testing.T) {
		serviceDir := filepath.Join(specsDir, "service")
		f1, err1 := ssacparser.ParseDir(serviceDir)
		f2, err2 := ssacparser.ParseDir(serviceDir)
		if err1 != nil || err2 != nil {
			t.Fatalf("parse errors: %v / %v", err1, err2)
		}
		if !reflect.DeepEqual(f1, f2) {
			t.Errorf("SSaC: two parses differ — len %d vs %d", len(f1), len(f2))
			for i := range f1 {
				if i >= len(f2) {
					break
				}
				if !reflect.DeepEqual(f1[i], f2[i]) {
					t.Errorf("  diff at [%d]: %s vs %s", i, f1[i].Name, f2[i].Name)
				}
			}
		}
	})

	t.Run("STML", func(t *testing.T) {
		frontendDir := filepath.Join(specsDir, "frontend")
		if _, err := os.Stat(frontendDir); err != nil {
			t.Skip("no frontend/ dir")
		}
		p1, err1 := stmlparser.ParseDir(frontendDir)
		p2, err2 := stmlparser.ParseDir(frontendDir)
		if err1 != nil || err2 != nil {
			t.Fatalf("parse errors: %v / %v", err1, err2)
		}
		if !reflect.DeepEqual(p1, p2) {
			t.Errorf("STML: two parses differ — len %d vs %d", len(p1), len(p2))
		}
	})

	t.Run("States", func(t *testing.T) {
		statesDir := filepath.Join(specsDir, "states")
		if _, err := os.Stat(statesDir); err != nil {
			t.Skip("no states/ dir")
		}
		d1, err1 := statemachine.ParseDir(statesDir)
		d2, err2 := statemachine.ParseDir(statesDir)
		if err1 != nil || err2 != nil {
			t.Fatalf("parse errors: %v / %v", err1, err2)
		}
		if !reflect.DeepEqual(d1, d2) {
			t.Errorf("States: two parses differ — len %d vs %d", len(d1), len(d2))
		}
	})

	t.Run("Policy", func(t *testing.T) {
		policyDir := filepath.Join(specsDir, "policy")
		if _, err := os.Stat(policyDir); err != nil {
			t.Skip("no policy/ dir")
		}
		p1, err1 := policy.ParseDir(policyDir)
		p2, err2 := policy.ParseDir(policyDir)
		if err1 != nil || err2 != nil {
			t.Fatalf("parse errors: %v / %v", err1, err2)
		}
		if !reflect.DeepEqual(p1, p2) {
			t.Errorf("Policy: two parses differ — len %d vs %d", len(p1), len(p2))
		}
	})

	t.Run("FuncSpec", func(t *testing.T) {
		funcDir := filepath.Join(specsDir, "func")
		if _, err := os.Stat(funcDir); err != nil {
			t.Skip("no func/ dir")
		}
		s1, err1 := funcspec.ParseDir(funcDir)
		s2, err2 := funcspec.ParseDir(funcDir)
		if err1 != nil || err2 != nil {
			t.Fatalf("parse errors: %v / %v", err1, err2)
		}
		if !reflect.DeepEqual(s1, s2) {
			t.Errorf("FuncSpec: two parses differ — len %d vs %d", len(s1), len(s2))
		}
	})

	t.Run("ProjectConfig", func(t *testing.T) {
		c1, err1 := projectconfig.Load(specsDir)
		c2, err2 := projectconfig.Load(specsDir)
		if err1 != nil || err2 != nil {
			t.Fatalf("load errors: %v / %v", err1, err2)
		}
		if !reflect.DeepEqual(c1, c2) {
			t.Error("ProjectConfig: two parses produced different results")
		}
	})
}

// diffSymbolTables reports which fields differ between two SymbolTables.
func diffSymbolTables(t *testing.T, a, b *ssacvalidator.SymbolTable) {
	t.Helper()
	t.Error("SymbolTable: two parses produced different results")

	if !reflect.DeepEqual(a.DDLTables, b.DDLTables) {
		t.Error("  DDLTables differ")
		for k := range a.DDLTables {
			if !reflect.DeepEqual(a.DDLTables[k], b.DDLTables[k]) {
				t.Errorf("    table %q differs", k)
			}
		}
		for k := range b.DDLTables {
			if _, ok := a.DDLTables[k]; !ok {
				t.Errorf("    table %q only in second parse", k)
			}
		}
	}
	if !reflect.DeepEqual(a.Models, b.Models) {
		t.Error("  Models differ")
	}
	if !reflect.DeepEqual(a.Operations, b.Operations) {
		t.Error("  Operations differ")
	}
	if !reflect.DeepEqual(a.DTOs, b.DTOs) {
		t.Error("  DTOs differ")
	}
	if !reflect.DeepEqual(a.Funcs, b.Funcs) {
		t.Error("  Funcs differ")
	}
}
