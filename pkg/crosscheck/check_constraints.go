//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkConstraints — DDL VARCHAR/CHECK ↔ OpenAPI maxLength/enum/format (X-65~X-72)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
)

func checkConstraints(fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil || len(fs.DDLTables) == 0 {
		return nil
	}
	var errs []CrossError
	for _, t := range fs.DDLTables {
		errs = append(errs, checkTableConstraints(t, fs)...)
	}
	return errs
}
