//ff:func feature=crosscheck type=loader control=iteration dimension=1
//ff:what populateFunc — FuncSpec에서 함수명, Request 필드 추출
package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateFunc(g *rule.Ground, fs *fullend.Fullstack) {
	specs := make(rule.StringSet)
	allSpecs := append(fs.ProjectFuncSpecs, fs.FullendPkgSpecs...)
	for _, sp := range allSpecs {
		key := strings.ToLower(sp.Package + "." + sp.Name)
		specs[key] = true
		var reqFields []string
		for _, f := range sp.RequestFields {
			reqFields = append(reqFields, f.Name)
			g.Types["Func.request."+sp.Name+"."+f.Name] = f.Type
		}
		g.Schemas["Func.request."+sp.Name] = reqFields
	}
	g.Lookup["Func.spec"] = specs
}
