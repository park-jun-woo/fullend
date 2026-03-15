//ff:func feature=gen-hurl type=util
//ff:what $ref를 실제 스키마로 해석한다
package hurl

import "github.com/getkin/kin-openapi/openapi3"

// resolveSchema follows $ref to get the actual schema.
func resolveSchema(ref *openapi3.SchemaRef) *openapi3.Schema {
	if ref == nil {
		return nil
	}
	return ref.Value
}
