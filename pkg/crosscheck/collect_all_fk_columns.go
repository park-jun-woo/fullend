//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectAllFKColumns — 모든 DDL 테이블의 FK 컬럼명 수집
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func collectAllFKColumns(fs *fullend.Fullstack) rule.StringSet {
	fks := make(rule.StringSet)
	for _, t := range fs.DDLTables {
		for _, fk := range t.ForeignKeys {
			fks[fk.Column] = true
		}
	}
	return fks
}
