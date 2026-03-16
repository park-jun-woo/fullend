//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what @call func에서 I/O 패키지 import 금지 검증
package crosscheck

import (
	"fmt"

	"github.com/geul-org/fullend/internal/funcspec"
)

// validateCallImports checks that a func spec does not import forbidden I/O packages.
func validateCallImports(ctx string, spec *funcspec.FuncSpec) []CrossError {
	var errs []CrossError
	ioImports := checkForbiddenImports(spec.Imports)
	for _, imp := range ioImports {
		errs = append(errs, CrossError{
			Rule:    "Func ↔ SSaC",
			Context: ctx,
			Message: fmt.Sprintf("@call func에서 I/O 패키지 %q import 금지. @call func은 순수 계산/판단 로직만 허용됩니다. DB, 네트워크, 파일 등 I/O가 필요하면 @model을 활용하세요.", imp),
			Level:   "ERROR",
		})
	}
	return errs
}
