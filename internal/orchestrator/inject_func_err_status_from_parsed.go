//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what injectFuncErrStatusFromParsed injects @error annotations from func specs into symbol table.

package orchestrator

import (
	"strings"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	"github.com/park-jun-woo/fullend/internal/genapi"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// injectFuncErrStatusFromParsed uses pre-parsed func specs to inject @error annotations.
func injectFuncErrStatusFromParsed(st *ssacvalidator.SymbolTable, parsed *genapi.ParsedSSOTs) {
	var allSpecs []funcspec.FuncSpec
	allSpecs = append(allSpecs, parsed.FullendPkgSpecs...)
	allSpecs = append(allSpecs, parsed.ProjectFuncSpecs...)

	for _, fs := range allSpecs {
		if fs.ErrStatus == 0 || fs.Package == "" {
			continue
		}
		// Register as "pkg._func" model key (matches SSaC convention).
		modelKey := fs.Package + "._func"
		ms, exists := st.Models[modelKey]
		if !exists {
			ms = ssacvalidator.ModelSymbol{Methods: make(map[string]ssacvalidator.MethodInfo)}
		}
		funcName := strings.ToUpper(fs.Name[:1]) + fs.Name[1:]
		mi := ms.Methods[funcName]
		mi.ErrStatus = fs.ErrStatus
		ms.Methods[funcName] = mi
		st.Models[modelKey] = ms
	}
}
