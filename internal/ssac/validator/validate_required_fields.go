//ff:func feature=ssac-validate type=rule control=iteration dimension=1
//ff:what 타입별 필수 필드 누락 검증
package validator

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validateRequiredFields는 타입별 필수 필드 누락을 검증한다.
func validateRequiredFields(sf parser.ServiceFunc) []ValidationError {
	var errs []ValidationError

	for i, seq := range sf.Sequences {
		ctx := errCtx{sf.FileName, sf.Name, i}

		errs = append(errs, validateSeqRequiredFields(seq, ctx)...)

		// Model 형식 검증: "Model.Method" 또는 "pkg.Func"
		if seq.Model == "" {
			continue
		}
		parts := strings.SplitN(seq.Model, ".", 2)
		if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
			errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("\"Model.Method\" 형식이어야 함: %q", seq.Model)))
		}
	}

	return errs
}
