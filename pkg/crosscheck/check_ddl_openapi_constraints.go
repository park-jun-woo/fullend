//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkDDLOpenAPIConstraints — DDL VARCHAR/CHECK ↔ OpenAPI maxLength/enum (X-66~X-70)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkDDLOpenAPIConstraints(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil || len(fs.DDLTables) == 0 {
		return nil
	}
	var errs []CrossError
	for _, t := range fs.DDLTables {
		errs = append(errs, checkDDLCheckEnum(g, t.Name, t.CheckEnums)...)
	}
	return errs
}
