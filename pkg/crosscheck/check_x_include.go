//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkXInclude — x-include 형식 검증 + DDL table/FK 존재 검증 (X-4, X-5, X-6)
package crosscheck

import (

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkXInclude(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil || len(fs.DDLTables) == 0 {
		return nil
	}
	var errs []CrossError
	for path, item := range fs.OpenAPIDoc.Paths.Map() {
		errs = append(errs, checkPathXInclude(g, path, item.Operations())...)
	}
	return errs
}
