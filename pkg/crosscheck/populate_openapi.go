//ff:func feature=crosscheck type=loader control=iteration dimension=1
//ff:what populateOpenAPI — OpenAPI에서 operationId, path, method, response 스키마 추출
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateOpenAPI(g *rule.Ground, fs *fullend.Fullstack) {
	if fs.OpenAPIDoc == nil {
		return
	}
	opIDs := make(rule.StringSet)
	paths := make(rule.StringSet)
	security := make(rule.StringSet)

	for path, item := range fs.OpenAPIDoc.Paths.Map() {
		paths[path] = true
		methods := make(rule.StringSet)
		populatePathOps(g, opIDs, methods, item.Operations())
		g.Lookup["OpenAPI.method."+path] = methods
	}
	g.Lookup["OpenAPI.operationId"] = opIDs
	g.Lookup["OpenAPI.path"] = paths

	if fs.OpenAPIDoc.Components != nil {
		for name := range fs.OpenAPIDoc.Components.SecuritySchemes {
			security[name] = true
		}
	}
	g.Lookup["OpenAPI.security"] = security
}
