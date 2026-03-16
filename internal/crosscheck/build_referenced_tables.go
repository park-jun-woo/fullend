//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what SSaC @model/@result에서 참조하는 DDL 테이블 이름을 수집
package crosscheck

import (
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// buildReferencedTables collects DDL table names referenced by SSaC @model and @result.
func buildReferencedTables(funcs []ssacparser.ServiceFunc) map[string]bool {
	tables := make(map[string]bool)
	for _, fn := range funcs {
		for _, seq := range fn.Sequences {
			collectReferencedTable(seq, tables)
		}
	}
	return tables
}
