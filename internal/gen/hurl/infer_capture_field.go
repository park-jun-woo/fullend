//ff:func feature=gen-hurl type=util control=iteration dimension=1
//ff:what 응답에서 캡처할 ID 필드를 추론한다
package hurl

import (
	"github.com/ettle/strcase"
	"github.com/getkin/kin-openapi/openapi3"
)

// inferCaptureField finds the response field that contains an ID to capture.
// e.g. response has "gig" object with "id" → capture "gig_id" from "$.gig.id"
// Checks both "id" (OpenAPI convention) and "ID" (Go convention) property names.
func inferCaptureField(respSchema *openapi3.Schema) (varName, jsonPath string) {
	if respSchema == nil {
		return "", ""
	}
	for name, propRef := range respSchema.Properties {
		prop := propRef.Value
		if prop.Type.Slice()[0] != "object" || prop.Properties == nil {
			continue
		}
		if _, hasID := prop.Properties["id"]; hasID {
			return strcase.ToSnake(name) + "_id", "$." + name + ".id"
		}
		if _, hasID := prop.Properties["ID"]; hasID {
			return strcase.ToSnake(name) + "_id", "$." + name + ".ID"
		}
	}
	return "", ""
}
