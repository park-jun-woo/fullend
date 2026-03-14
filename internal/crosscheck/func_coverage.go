package crosscheck

import (
	"fmt"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// CheckFuncCoverage warns about project func specs not referenced by any SSaC @call.
func CheckFuncCoverage(
	funcs []ssacparser.ServiceFunc,
	projectFuncSpecs []funcspec.FuncSpec,
) []CrossError {
	// Collect pkg.Function names referenced by SSaC @call sequences.
	referenced := make(map[string]bool)
	for _, fn := range funcs {
		for _, seq := range fn.Sequences {
			if seq.Type == "call" && seq.Model != "" {
				referenced[seq.Model] = true
			}
		}
	}

	var errs []CrossError
	for _, spec := range projectFuncSpecs {
		key := spec.Package + "." + spec.Name
		if !referenced[key] {
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
