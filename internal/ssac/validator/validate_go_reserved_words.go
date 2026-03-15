//ff:func feature=ssac-validate type=rule
//ff:what Go 예약어 충돌 검증

package validator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// goReservedWords는 Go 예약어 25개다.
var goReservedWords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true,
	"continue": true, "default": true, "defer": true, "else": true,
	"fallthrough": true, "for": true, "func": true, "go": true,
	"goto": true, "if": true, "import": true, "interface": true,
	"map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true,
	"var": true,
}

// validateGoReservedWords는 SSaC Inputs 키가 Go 예약어와 충돌하면 ERROR를 반환한다.
func validateGoReservedWords(funcs []parser.ServiceFunc, st *SymbolTable) []ValidationError {
	var errs []ValidationError
	seen := map[string]bool{} // 중복 에러 방지: "table.column"

	for _, sf := range funcs {
		for i, seq := range sf.Sequences {
			if seq.Package != "" || seq.Type == parser.SeqCall {
				continue // 패키지 모델과 @call은 models_gen.go 대상 아님
			}
			for key := range seq.Inputs {
				paramName := toLowerFirst(key)
				if !goReservedWords[paramName] {
					continue
				}
				// DDL 테이블에서 컬럼 역추적
				snakeName := toSnakeCase(key)
				tableName, found := findColumnTable(snakeName, seq.Model, st)
				ctx := errCtx{sf.FileName, sf.Name, i}
				dedup := tableName + "." + snakeName
				if seen[dedup] {
					continue
				}
				seen[dedup] = true
				if found {
					errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("DDL column %q in table %q is a Go reserved word — rename the column (e.g. %q)", snakeName, tableName, "tx_"+snakeName)))
				} else {
					errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("parameter name %q is a Go reserved word — rename the DDL column", paramName)))
				}
			}
		}
	}
	return errs
}
