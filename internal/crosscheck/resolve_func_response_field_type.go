//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-openapi
//ff:what func spec Response 필드에서 필드 타입 조회
package crosscheck

import "github.com/geul-org/fullend/internal/funcspec"

// resolveFuncResponseFieldType looks up a field type from func spec Response fields.
func resolveFuncResponseFieldType(specs []funcspec.FuncSpec, respTypeName, field string) string {
	for _, spec := range specs {
		if spec.Name+"Response" != respTypeName {
			continue
		}
		return findFieldType(spec.ResponseFields, field)
	}
	return ""
}
