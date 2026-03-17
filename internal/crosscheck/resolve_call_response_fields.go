//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-openapi
//ff:what @call 시퀀스의 func spec에서 응답 필드 해석
package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// resolveCallResponseFields resolves response fields from a @call sequence's func spec.
func resolveCallResponseFields(seq ssacparser.Sequence, funcSpecs []funcspec.FuncSpec) []string {
	typeName := seq.Result.Type
	if idx := strings.LastIndex(typeName, "."); idx >= 0 {
		typeName = typeName[idx+1:]
	}
	for _, fs := range funcSpecs {
		if ucFirst(fs.Name)+"Response" == typeName {
			return extractFuncSpecFieldKeys(fs)
		}
	}
	return nil
}
