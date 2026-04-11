//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkVarcharMaxLength — DDL VARCHAR(n) ↔ OpenAPI maxLength (X-67, X-70)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
)

func checkVarcharMaxLength(fs *fullend.Fullstack) []CrossError {
	if len(fs.DDLTables) == 0 || len(fs.ResponseConstraints) == 0 {
		return nil
	}
	var errs []CrossError
	for _, t := range fs.DDLTables {
		for col, vLen := range t.VarcharLen {
			errs = append(errs, checkColMaxLength(t.Name, col, vLen, fs)...)
		}
	}
	return errs
}
