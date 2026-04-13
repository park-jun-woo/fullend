//ff:func feature=rule type=util control=iteration dimension=1
//ff:what extractPathParams — openapi3.Parameters 에서 path 파라미터 추출
package ground

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func extractPathParams(params openapi3.Parameters) []rule.PathParam {
	var result []rule.PathParam
	for _, p := range params {
		if p.Value == nil || p.Value.In != "path" {
			continue
		}
		result = append(result, rule.PathParam{
			Name:   p.Value.Name,
			GoType: schemaGoType(p.Value.Schema),
		})
	}
	return result
}

func schemaGoType(ref *openapi3.SchemaRef) string {
	if ref == nil || ref.Value == nil {
		return "string"
	}
	t := ref.Value.Type
	switch {
	case t.Is("integer"):
		if ref.Value.Format == "int32" {
			return "int32"
		}
		return "int64"
	case t.Is("number"):
		return "float64"
	case t.Is("boolean"):
		return "bool"
	default:
		return "string"
	}
}
