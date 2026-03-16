//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what 응답 참조에서 JSON 스키마 속성명을 맵에 추가
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// addResponseSchemaProps adds schema property names from a response reference.
func addResponseSchemaProps(props map[string]bool, respRef *openapi3.ResponseRef) {
	if respRef == nil || respRef.Value == nil || respRef.Value.Content == nil {
		return
	}
	ct := respRef.Value.Content.Get("application/json")
	if ct == nil || ct.Schema == nil {
		return
	}
	schema := resolveSchemaRef(ct.Schema)
	if schema == nil {
		return
	}
	for propName := range schema.Properties {
		props[propName] = true
	}
}
