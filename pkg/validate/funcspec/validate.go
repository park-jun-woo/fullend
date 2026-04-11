//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what Validate — func spec built-in 패키지명 충돌 검증 (F-1)
package funcspec

import (
	parsefunc "github.com/park-jun-woo/fullend/pkg/parser/funcspec"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

var builtinPackages = map[string]bool{
	"auth": true, "session": true, "cache": true, "file": true,
	"queue": true, "crypto": true, "storage": true, "mail": true,
	"text": true, "image": true,
}

// Validate checks if project func specs conflict with built-in package names.
func Validate(projectSpecs []parsefunc.FuncSpec, builtinSpecs []parsefunc.FuncSpec) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, sp := range projectSpecs {
		if builtinPackages[sp.Package] {
			errs = append(errs, validate.ValidationError{
				Rule: "F-1", File: sp.Package + "/" + sp.Name + ".go", Level: "WARNING",
				Message: "func package " + sp.Package + " overrides built-in package", SeqIdx: -1,
			})
		}
	}
	return errs
}
