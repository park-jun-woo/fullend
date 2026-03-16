//ff:func feature=ssac-validate type=rule control=iteration dimension=2
//ff:what SSaC query 사용과 OpenAPI x-extensions 교차 검증

package validator

import (
	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validateQueryUsage는 SSaC의 query 사용과 OpenAPI x-extensions 간 교차 검증을 수행한다.
func validateQueryUsage(sf parser.ServiceFunc, st *SymbolTable) []ValidationError {
	if st == nil {
		return nil
	}

	op, hasOp := st.Operations[sf.Name]
	opHasQueryOpts := hasOp && op.HasQueryOpts()

	specHasQuery := false
	for _, seq := range sf.Sequences {
		if specHasQuery {
			break
		}
		for _, val := range seq.Inputs {
			if val != "query" {
				continue
			}
			specHasQuery = true
			break
		}
	}

	var errs []ValidationError
	ctx := errCtx{sf.FileName, sf.Name, -1}

	if specHasQuery && !opHasQueryOpts {
		errs = append(errs, ctx.err("@query", "SSaC에 query가 사용되었지만 OpenAPI에 x-pagination/sort/filter가 없습니다"))
	}
	if opHasQueryOpts && !specHasQuery {
		errs = append(errs, ValidationError{
			FileName: ctx.fileName, FuncName: ctx.funcName, SeqIndex: ctx.seqIndex,
			Tag: "@query", Message: "OpenAPI에 x-pagination/sort/filter가 있지만 SSaC에 query가 사용되지 않았습니다", Level: "WARNING",
		})
	}

	return errs
}
