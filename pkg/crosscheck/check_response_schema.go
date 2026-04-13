//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkResponseSchema — @response fields ↔ OpenAPI response 스키마 (X-17, X-18)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkResponseSchema(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil {
		return nil
	}
	var errs []CrossError
	for funcName, fields := range g.Schemas {
		if len(funcName) < 14 || funcName[:14] != "SSaC.response." {
			continue
		}
		opID := funcName[14:]
		schemaKey := "OpenAPI.response." + opID
		target, ok := g.Schemas[schemaKey]
		if !ok {
			continue
		}
		errs = append(errs, evalSchemaMatch(g, fields, target, "X-17", opID)...)
		errs = append(errs, evalReverseCoverage(target, fields, "X-18", opID)...)
	}
	return errs
}
