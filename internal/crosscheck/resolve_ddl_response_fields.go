//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what DDL 컬럼에서 모델 시퀀스의 응답 필드 해석
package crosscheck

import (
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// resolveDDLResponseFields resolves response fields from DDL columns for a model sequence.
func resolveDDLResponseFields(seq ssacparser.Sequence, st *ssacvalidator.SymbolTable) []string {
	if st == nil {
		return nil
	}
	tableName := seq.Result.Type
	for tbl, ddl := range st.DDLTables {
		if !matchTableType(tbl, tableName) {
			continue
		}
		return sortedColumnNames(ddl.Columns)
	}
	return nil
}
