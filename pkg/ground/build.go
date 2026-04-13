//ff:func feature=rule type=loader control=sequence
//ff:what Build — Fullstack에서 완전한 rule.Ground를 구축
package ground

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// Build extracts all lookup data from Fullstack into a shared rule.Ground.
// Used by pkg/validate, pkg/crosscheck, and pkg/generate.
func Build(fs *fullend.Fullstack) *rule.Ground {
	g := &rule.Ground{
		Lookup:     make(map[string]rule.StringSet),
		Types:      make(map[string]string),
		Pairs:      make(map[string]rule.StringSet),
		Config:     make(map[string]bool),
		Vars:       make(rule.StringSet),
		Flags:      make(rule.StringSet),
		Schemas:    make(map[string][]string),
		Models:     make(map[string]rule.ModelInfo),
		Tables:     make(map[string]rule.TableInfo),
		Ops:        make(map[string]rule.OperationInfo),
		ReqSchemas: make(map[string]rule.RequestSchemaInfo),
	}
	// 기존 populate — validate/crosscheck 소비
	populateOpenAPI(g, fs)
	populateSSaC(g, fs)
	populateStates(g, fs)
	populateFunc(g, fs)
	populateManifest(g, fs)
	populateDDL(g, fs)
	populateRego(g, fs)
	populateOpenAPIConstraints(g, fs)
	populateOpenAPIParams(g, fs)
	populateVarTypes(g, fs)
	populateGoReservedWords(g)
	populateHurl(g, fs)

	// 신규 populate — generate 소비 (Phase002)
	populateModels(g, fs)
	populateTables(g, fs)
	populateOps(g, fs)
	populateRequestSchemas(g, fs)

	// legacy bridge — validate 규칙(S-48 등) 이 아직 Lookup["SymbolTable.model"] 을 쓰므로
	// g.Models 에서 파생해 유지. Phase003 에서 validate 마이그 후 제거 예정.
	populateModelLookup(g)

	return g
}
