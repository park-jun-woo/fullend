//ff:func feature=crosscheck type=loader control=sequence
//ff:what BuildGround — Fullstack에서 rule.Ground를 구축
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// BuildGround extracts lookup data from Fullstack into a rule.Ground.
func BuildGround(fs *fullend.Fullstack) *rule.Ground {
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
	populateHurl(g, fs)
	return g
}
