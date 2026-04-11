//ff:func feature=rule type=loader control=iteration dimension=1
//ff:what populateOpenAPIParams — operationId별 path 파라미터, x-sort, x-filter 등록
package ground

import (

	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateOpenAPIParams(g *rule.Ground, fs *fullend.Fullstack) {
	if fs.OpenAPIDoc == nil {
		return
	}
	for _, item := range fs.OpenAPIDoc.Paths.Map() {
		populatePathOpsParams(g, item.Operations())
	}
}
