//ff:func feature=gen-hurl type=util control=sequence
//ff:what Checks if the OpenAPI doc has any security schemes defined.
package hurl

import "github.com/getkin/kin-openapi/openapi3"

func hasSecuritySchemes(doc *openapi3.T) bool {
	if doc.Components == nil {
		return false
	}
	return len(doc.Components.SecuritySchemes) > 0
}
