//ff:func feature=rule type=loader control=sequence
//ff:what populateModels 검증 — iface + sqlc + FuncSpec 결합
package ground

import (
	"testing"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/parser/funcspec"
	"github.com/park-jun-woo/fullend/pkg/parser/iface"
	"github.com/park-jun-woo/fullend/pkg/parser/sqlc"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func TestPopulateModelsFromInterfaces(t *testing.T) {
	g := newGround()
	fs := &fullend.Fullstack{
		ModelInterfaces: []iface.Interface{
			{Name: "UserModel", Methods: []string{"Create", "FindByID"}},
		},
	}
	populateModels(g, fs)

	m, ok := g.Models["UserModel"]
	if !ok {
		t.Fatal("UserModel not populated")
	}
	if len(m.Methods) != 2 {
		t.Errorf("expected 2 methods, got %d", len(m.Methods))
	}
	if _, ok := m.Methods["Create"]; !ok {
		t.Error("Create method missing")
	}
}

func TestPopulateModelsSqlcMergesCardinality(t *testing.T) {
	g := newGround()
	fs := &fullend.Fullstack{
		SqlcQueries: []sqlc.Query{
			{Model: "User", Name: "FindByID", Cardinality: "one", Params: []string{"ID"}},
			{Model: "User", Name: "List", Cardinality: "many"},
		},
	}
	populateModels(g, fs)

	m, ok := g.Models["User"]
	if !ok {
		t.Fatal("User not populated")
	}
	if m.Methods["FindByID"].Cardinality != "one" {
		t.Errorf("FindByID cardinality: got %q", m.Methods["FindByID"].Cardinality)
	}
	if len(m.Methods["FindByID"].Params) != 1 || m.Methods["FindByID"].Params[0] != "ID" {
		t.Errorf("FindByID params: got %v", m.Methods["FindByID"].Params)
	}
	if m.Methods["List"].Cardinality != "many" {
		t.Errorf("List cardinality: got %q", m.Methods["List"].Cardinality)
	}
}

func TestPopulateModelsInjectsErrStatus(t *testing.T) {
	g := newGround()
	fs := &fullend.Fullstack{
		ProjectFuncSpecs: []funcspec.FuncSpec{
			{Package: "auth", Name: "hashPassword", ErrStatus: 500},
			{Package: "billing", Name: "holdEscrow", ErrStatus: 402},
		},
	}
	populateModels(g, fs)

	auth, ok := g.Models["auth._func"]
	if !ok {
		t.Fatal("auth._func not populated")
	}
	if auth.Methods["HashPassword"].ErrStatus != 500 {
		t.Errorf("auth.HashPassword.ErrStatus: got %d", auth.Methods["HashPassword"].ErrStatus)
	}
	billing, ok := g.Models["billing._func"]
	if !ok || billing.Methods["HoldEscrow"].ErrStatus != 402 {
		t.Errorf("billing.HoldEscrow.ErrStatus wrong: %v", billing)
	}
}

func TestPopulateModelsErrStatusSkipsZero(t *testing.T) {
	g := newGround()
	fs := &fullend.Fullstack{
		ProjectFuncSpecs: []funcspec.FuncSpec{
			{Package: "auth", Name: "hashPassword", ErrStatus: 0}, // 미지정 → 스킵
		},
	}
	populateModels(g, fs)

	if _, ok := g.Models["auth._func"]; ok {
		t.Error("auth._func should not be created when ErrStatus == 0")
	}
}

func newGround() *rule.Ground {
	return &rule.Ground{
		Lookup:  make(map[string]rule.StringSet),
		Types:   make(map[string]string),
		Pairs:   make(map[string]rule.StringSet),
		Config:  make(map[string]bool),
		Vars:    make(rule.StringSet),
		Flags:   make(rule.StringSet),
		Schemas: make(map[string][]string),
		Models:  make(map[string]rule.ModelInfo),
	}
}
