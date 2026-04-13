//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkFuncPurity — func body TODO 감지 (X-40) + 금지 import 감지 (X-41)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
)

func checkFuncPurity(fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	for _, sp := range fs.ProjectFuncSpecs {
		errs = append(errs, checkSingleFuncPurity(sp.Package, sp.Name, sp.HasBody, sp.Imports)...)
	}
	return errs
}
