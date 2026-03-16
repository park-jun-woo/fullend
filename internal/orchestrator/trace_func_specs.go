//ff:func feature=orchestrator type=util control=iteration dimension=2
//ff:what traceFuncSpecs finds func specs referenced by @call sequences.

package orchestrator

import (
	"strings"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

func traceFuncSpecs(sf *ssacparser.ServiceFunc, specs []funcspec.FuncSpec, specsDir string) []ChainLink {
	callPkgFuncs := map[string]string{} // "pkg.Func" -> pkg
	for _, seq := range sf.Sequences {
		if seq.Type != "call" || seq.Model == "" {
			continue
		}
		parts := strings.SplitN(seq.Model, ".", 2)
		if len(parts) == 2 {
			callPkgFuncs[seq.Model] = parts[0]
		}
	}

	if len(callPkgFuncs) == 0 {
		return nil
	}

	var links []ChainLink
	for callRef, pkg := range callPkgFuncs {
		parts := strings.SplitN(callRef, ".", 2)
		funcName := ""
		if len(parts) == 2 {
			funcName = parts[1]
		}
		if link, ok := findFuncSpecLink(callRef, pkg, funcName, specs, specsDir); ok {
			links = append(links, link)
		}
	}
	return links
}
