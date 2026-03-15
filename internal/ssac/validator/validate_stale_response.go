//ff:func feature=ssac-validate type=rule
//ff:what put/delete 이후 갱신 없이 response에서 사용되는 변수 경고

package validator

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validateStaleResponse는 put/delete 이후 갱신 없이 response에서 사용되는 변수를 경고한다.
func validateStaleResponse(sf parser.ServiceFunc) []ValidationError {
	var errs []ValidationError

	getVars := map[string]string{}   // var → model
	mutated := map[string]bool{}     // model → mutated?

	for i, seq := range sf.Sequences {
		switch seq.Type {
		case parser.SeqGet:
			if seq.Result != nil && seq.Model != "" {
				modelName := strings.SplitN(seq.Model, ".", 2)[0]
				getVars[seq.Result.Var] = modelName
				mutated[modelName] = false
			}
		case parser.SeqPut, parser.SeqDelete:
			if seq.Model != "" {
				modelName := strings.SplitN(seq.Model, ".", 2)[0]
				mutated[modelName] = true
			}
		case parser.SeqResponse:
			if seq.SuppressWarn {
				continue
			}
			ctx := errCtx{sf.FileName, sf.Name, i}
			for field, val := range seq.Fields {
				ref := rootVar(val)
				if modelName, ok := getVars[ref]; ok && mutated[modelName] {
					errs = append(errs, ValidationError{
						FileName: ctx.fileName,
						FuncName: ctx.funcName,
						SeqIndex: ctx.seqIndex,
						Tag:      "@response",
						Message:  fmt.Sprintf("%q (필드 %q)가 %s 수정 이후 갱신 없이 response에 사용됩니다", ref, field, modelName),
						Level:    "WARNING",
					})
				}
			}
		}
	}

	return errs
}
