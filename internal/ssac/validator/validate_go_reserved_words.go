//ff:func feature=ssac-validate type=rule control=iteration dimension=2 topic=string-convert
//ff:what Go 예약어 충돌 검증

package validator

import "github.com/geul-org/fullend/internal/ssac/parser"

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
			ctx := errCtx{sf.FileName, sf.Name, i}
			errs = append(errs, checkReservedKeys(seq, ctx, st, seen)...)
		}
	}
	return errs
}
