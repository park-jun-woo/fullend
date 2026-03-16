//ff:func feature=ssac-validate type=rule control=iteration dimension=1
//ff:what 변수 선언 전 참조 검증

package validator

import "github.com/geul-org/fullend/internal/ssac/parser"

// validateVariableFlow는 변수가 선언 전에 참조되지 않는지 검증한다.
func validateVariableFlow(sf parser.ServiceFunc) []ValidationError {
	var errs []ValidationError
	declared := map[string]bool{
		"currentUser": true,
	}
	if sf.Subscribe != nil {
		declared["message"] = true
	}

	for i, seq := range sf.Sequences {
		ctx := errCtx{sf.FileName, sf.Name, i}
		errs = checkTargetDeclared(seq, ctx, declared, errs)
		errs = checkArgsDeclared(seq, ctx, declared, errs)
		errs = checkInputsDeclared(seq, ctx, declared, errs)
		errs = checkFieldsDeclared(seq, ctx, declared, errs)

		// Result로 변수 선언
		if seq.Result != nil {
			declared[seq.Result.Var] = true
		}
	}

	return errs
}
