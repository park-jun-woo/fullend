//ff:func feature=crosscheck type=util control=sequence topic=config-check
//ff:what OpenAPI securitySchemes 이름을 수집
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

func collectSchemeNames(doc *openapi3.T) map[string]bool {
	schemeNames := make(map[string]bool)
	if doc.Components != nil && doc.Components.SecuritySchemes != nil {
		for name := range doc.Components.SecuritySchemes {
			schemeNames[name] = true
		}
	}
	return schemeNames
}
