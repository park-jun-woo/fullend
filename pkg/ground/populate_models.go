//ff:func feature=rule type=loader control=sequence
//ff:what populateModels — iface + sqlc + FuncSpec 를 결합해 g.Models 를 채움
package ground

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/parser/funcspec"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// populateModels builds g.Models by merging three sources:
//  1. ModelInterfaces (Go iface) — initial model / method skeleton
//  2. SqlcQueries — method cardinality + params
//  3. FuncSpecs (FullendPkg + Project) — @error → ErrStatus on "<pkg>._func" models
func populateModels(g *rule.Ground, fs *fullend.Fullstack) {
	if g.Models == nil {
		g.Models = make(map[string]rule.ModelInfo)
	}
	seedFromInterfaces(g, fs.ModelInterfaces)
	mergeSqlcQueries(g, fs.SqlcQueries)
	injectErrStatus(g, fs.FullendPkgSpecs)
	injectErrStatus(g, fs.ProjectFuncSpecs)
}

func injectErrStatus(g *rule.Ground, specs []funcspec.FuncSpec) {
	for _, spec := range specs {
		if spec.ErrStatus == 0 || spec.Package == "" {
			continue
		}
		modelKey := spec.Package + "._func"
		info := ensureModel(g, modelKey)
		funcName := upperFirst(spec.Name)
		mi := info.Methods[funcName]
		mi.ErrStatus = spec.ErrStatus
		info.Methods[funcName] = mi
		g.Models[modelKey] = info
	}
}
