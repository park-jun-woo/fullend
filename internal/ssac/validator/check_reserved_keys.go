//ff:func feature=ssac-validate type=util control=iteration dimension=1 topic=string-convert
//ff:what 단일 시퀀스의 Inputs 키에서 Go 예약어 충돌을 검사한다
package validator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// checkReservedKeys는 단일 시퀀스의 Inputs 키에서 Go 예약어 충돌을 검사한다.
func checkReservedKeys(seq parser.Sequence, ctx errCtx, st *SymbolTable, seen map[string]bool) []ValidationError {
	var errs []ValidationError
	for key := range seq.Inputs {
		paramName := toLowerFirst(key)
		if !goReservedWords[paramName] {
			continue
		}
		snakeName := toSnakeCase(key)
		tableName, found := findColumnTable(snakeName, seq.Model, st)
		dedup := tableName + "." + snakeName
		if seen[dedup] {
			continue
		}
		seen[dedup] = true
		if !found {
			errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("parameter name %q is a Go reserved word — rename the DDL column", paramName)))
			continue
		}
		errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("DDL column %q in table %q is a Go reserved word — rename the column (e.g. %q)", snakeName, tableName, "tx_"+snakeName)))
	}
	return errs
}
