//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=config-check
//ff:what fullend.yaml 미들웨어와 OpenAPI securitySchemes 일치 여부 검증
package crosscheck

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// CheckMiddleware validates that fullend.yaml middleware matches OpenAPI securitySchemes.
func CheckMiddleware(middleware []string, doc *openapi3.T) []CrossError {
	var errs []CrossError

	schemeNames := collectSchemeNames(doc)

	mwSet := make(map[string]bool)
	for _, m := range middleware {
		mwSet[m] = true
	}

	// Rule 1: Each OpenAPI securityScheme must have a matching middleware.
	for name := range schemeNames {
		if !mwSet[name] {
			errs = append(errs, CrossError{
				Rule:       "Config ↔ OpenAPI",
				Context:    fmt.Sprintf("securitySchemes.%s", name),
				Message:    fmt.Sprintf("OpenAPI securityScheme %q has no matching middleware in fullend.yaml", name),
				Suggestion: fmt.Sprintf("fullend.yaml backend.middleware에 %q 추가", name),
			})
		}
	}

	// Rule 2: Each middleware must have a matching OpenAPI securityScheme.
	for _, m := range middleware {
		if !schemeNames[m] {
			errs = append(errs, CrossError{
				Rule:       "Config ↔ OpenAPI",
				Context:    fmt.Sprintf("middleware.%s", m),
				Message:    fmt.Sprintf("middleware %q has no matching OpenAPI securityScheme", m),
				Suggestion: fmt.Sprintf("OpenAPI components.securitySchemes에 %q 추가", m),
			})
		}
	}

	// Rule 3: Endpoints referencing security names not in middleware.
	errs = append(errs, checkEndpointSecurity(doc, mwSet)...)

	return errs
}
