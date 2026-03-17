//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=import-collect
//ff:what HTTP 함수에 필요한 import 경로를 수집
package generator

import (
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func collectImports(sf parser.ServiceFunc, reqParams []typedRequestParam, pathParams []validator.PathParam, needsCU bool, needsQO bool) []string {
	seen := map[string]bool{
		"net/http":                  true,
		"github.com/gin-gonic/gin": true,
	}

	collectSeqImports(sf, seen)
	collectParamTypeImports(reqParams, seen)
	collectPathParamImports(pathParams, seen)

	if needsCU || needsQO {
		seen["model"] = true
	}
	if hasWriteSequence(sf.Sequences) {
		seen["database/sql"] = true
	}

	imports := buildOrderedImports(seen)

	for _, imp := range sf.Imports {
		imports = append(imports, imp)
	}
	return imports
}
