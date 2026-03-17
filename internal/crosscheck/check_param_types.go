//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=ssac-ddl
//ff:what SSaC @param 입력이 DDL 컬럼에 존재하는지 검증
package crosscheck

import (
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func checkParamTypes(seq ssacparser.Sequence, st *ssacvalidator.SymbolTable, ctx string, seqIdx int) []CrossError {
	var errs []CrossError

	parts := strings.SplitN(seq.Model, ".", 2)
	if len(parts) < 2 {
		return errs
	}
	modelName := parts[0]
	tableName := modelToTable(modelName)

	table, ok := st.DDLTables[tableName]
	if !ok {
		return errs
	}

	for key, value := range seq.Inputs {
		errs = append(errs, checkParamColumn(key, value, modelName, tableName, table, ctx, seqIdx)...)
	}

	return errs
}
