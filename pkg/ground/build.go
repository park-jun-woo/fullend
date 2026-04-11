//ff:func feature=rule type=loader control=sequence
//ff:what Build — Fullstack에서 완전한 rule.Ground를 구축
package ground

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// Build extracts all lookup data from Fullstack into a shared rule.Ground.
// Used by both pkg/validate and pkg/crosscheck.
func Build(fs *fullend.Fullstack) *rule.Ground {
	g := &rule.Ground{
		Lookup:  make(map[string]rule.StringSet),
		Types:   make(map[string]string),
		Pairs:   make(map[string]rule.StringSet),
		Config:  make(map[string]bool),
		Vars:    make(rule.StringSet),
		Flags:   make(rule.StringSet),
		Schemas: make(map[string][]string),
	}
	populateOpenAPI(g, fs)
	populateSSaC(g, fs)
	populateStates(g, fs)
	populateFunc(g, fs)
	populateManifest(g, fs)
	populateDDL(g, fs)
	populateRego(g, fs)
	populateOpenAPIConstraints(g, fs)
	populateOpenAPIParams(g, fs)
	populateSymbolTable(g, fs)
	populateVarTypes(g, fs)
	populateGoReservedWords(g)
	populateHurl(g, fs)
	return g
}
