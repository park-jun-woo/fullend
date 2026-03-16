//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what 프로젝트 func spec이 SSaC @call에서 참조되는지 검증
package crosscheck

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// CheckFuncCoverage warns about project func specs not referenced by any SSaC @call.
func CheckFuncCoverage(
	funcs []ssacparser.ServiceFunc,
	projectFuncSpecs []funcspec.FuncSpec,
) []CrossError {
	referenced := buildCallReferences(funcs)

	var errs []CrossError
	for _, spec := range projectFuncSpecs {
		key := spec.Package + "." + spec.Name
		if !referenced[strings.ToLower(key)] {
			errs = append(errs, CrossError{
				Rule:       "Func → SSaC",
				Context:    key,
				Message:    fmt.Sprintf("func spec %q is not referenced by any SSaC @call", key),
				Level:      "WARNING",
				Suggestion: fmt.Sprintf("SSaC에서 @call %s를 추가하거나 func/%s를 제거하세요", key, spec.Package),
			})
		}
	}
	return errs
}
