//ff:func feature=orchestrator type=loader control=sequence
//ff:what buildSTMLGround — STML 검증용 Ground 구축 (OpenAPI operationId 등록)
package orchestrator

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func buildSTMLGround(fs *fullend.Fullstack) *rule.Ground {
	ground := &rule.Ground{
		Lookup:  make(map[string]rule.StringSet),
		Types:   make(map[string]string),
		Pairs:   make(map[string]rule.StringSet),
		Config:  make(map[string]bool),
		Vars:    make(rule.StringSet),
		Flags:   make(rule.StringSet),
		Schemas: make(map[string][]string),
	}
	if fs.OpenAPIDoc != nil {
		opIDs := make(rule.StringSet)
		for _, item := range fs.OpenAPIDoc.Paths.Map() {
			populateSTMLOps(opIDs, item.Operations())
		}
		ground.Lookup["OpenAPI.operationId"] = opIDs
	}
	return ground
}
