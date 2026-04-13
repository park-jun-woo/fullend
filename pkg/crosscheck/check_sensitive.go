//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkSensitive — 민감 패턴 컬럼 @sensitive 미선언 경고 (X-61)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
)

func checkSensitive(fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	for _, t := range fs.DDLTables {
		errs = append(errs, checkTableSensitiveCols(t.Name, t.Columns)...)
	}
	return errs
}
