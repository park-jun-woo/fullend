//ff:func feature=crosscheck type=util control=sequence
//ff:what OpenAPI SchemaRef에서 Schema 값 추출
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

func resolveSchemaRef(ref *openapi3.SchemaRef) *openapi3.Schema {
	if ref == nil {
		return nil
	}
	return ref.Value
}
