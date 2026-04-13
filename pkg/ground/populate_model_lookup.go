//ff:func feature=rule type=loader control=iteration dimension=1
//ff:what populateModelLookup — g.Models 에서 모델 이름 집합을 Lookup 에 복사 (legacy)
package ground

import "github.com/park-jun-woo/fullend/pkg/rule"

// populateModelLookup mirrors g.Models keys into g.Lookup["SymbolTable.model"].
// Legacy compatibility for validate rules that still consume the flat lookup
// (e.g. S-48 via pkg/validate/ssac/validate_model_refs.go). Remove when
// validate migrates to g.Models directly (Phase003).
func populateModelLookup(g *rule.Ground) {
	models := make(rule.StringSet, len(g.Models))
	for k := range g.Models {
		models[k] = true
	}
	g.Lookup["SymbolTable.model"] = models
}
