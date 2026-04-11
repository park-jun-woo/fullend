//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkXIncludeFK — x-include FKColumn이 DDL FK 제약으로 선언되었는지 WARNING (X-6)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkXIncludeFK(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil || len(fs.DDLTables) == 0 {
		return nil
	}
	fkCols := collectAllFKColumns(fs)
	var errs []CrossError
	for path, item := range fs.OpenAPIDoc.Paths.Map() {
		errs = append(errs, checkXIncludeFKOps(path, item, fkCols)...)
	}
	_ = g
	return errs
}
