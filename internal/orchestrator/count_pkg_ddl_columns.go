//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what countPkgDDLColumns — pkg/parser/ddl Table에서 총 컬럼 수 집계
package orchestrator

import "github.com/park-jun-woo/fullend/pkg/parser/fullend"

func countPkgDDLColumns(fs *fullend.Fullstack) int {
	cols := 0
	for _, t := range fs.DDLTables {
		cols += len(t.Columns)
	}
	return cols
}
