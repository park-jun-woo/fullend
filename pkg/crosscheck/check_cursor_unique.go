//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkCursorUnique — cursor sort default가 UNIQUE 컬럼인지 검증 (X-8)
package crosscheck

import (

	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkCursorUnique(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil || len(fs.DDLTables) == 0 {
		return nil
	}
	var errs []CrossError
	for path, item := range fs.OpenAPIDoc.Paths.Map() {
		for _, op := range item.Operations() {
			errs = append(errs, checkOpCursorUnique(g, op, path)...)
		}
	}
	return errs
}
