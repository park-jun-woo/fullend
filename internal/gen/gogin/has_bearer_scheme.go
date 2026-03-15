//ff:func feature=gen-gogin type=util
//ff:what checks if the OpenAPI doc has a bearerAuth security scheme

package gogin

import "github.com/getkin/kin-openapi/openapi3"

// hasBearerScheme checks if the OpenAPI doc has a bearerAuth security scheme.
func hasBearerScheme(doc *openapi3.T) bool {
	if doc == nil || doc.Components == nil || doc.Components.SecuritySchemes == nil {
		return false
	}
	for _, ref := range doc.Components.SecuritySchemes {
		if ref.Value != nil && ref.Value.Type == "http" && ref.Value.Scheme == "bearer" {
			return true
		}
	}
	return false
}
