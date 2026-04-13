//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkDDLOpenAPICoverage — DDL column → OpenAPI response 포함 여부 (X-10)
package crosscheck

import (
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkDDLOpenAPICoverage(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil || len(fs.DDLTables) == 0 {
		return nil
	}
	var errs []CrossError
	for _, t := range fs.DDLTables {
		model := inflection.Singular(strings.Title(t.Name))
		schemaKey := "OpenAPI.response.resolved.Get" + model
		target := g.Schemas[schemaKey]
		if len(target) == 0 {
			continue
		}
		errs = append(errs, evalDDLColumnCoverage(g, t.Name, target)...)
	}
	return errs
}
