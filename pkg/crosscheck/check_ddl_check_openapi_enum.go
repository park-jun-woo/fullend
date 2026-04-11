//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkDDLCheckOpenAPIEnum — DDL CHECK values ↔ OpenAPI enum 일치 검증 (X-69)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkDDLCheckOpenAPIEnum(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.DDLTables) == 0 || len(fs.ResponseConstraints) == 0 {
		return nil
	}
	var errs []CrossError
	for _, t := range fs.DDLTables {
		errs = append(errs, compareDDLEnumWithOpenAPICols(t.Name, t.CheckEnums, fs)...)
	}
	_ = g
	return errs
}
