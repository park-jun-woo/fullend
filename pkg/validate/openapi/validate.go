//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what Validate — OpenAPI path 파라미터 충돌 검증 (O-1)
package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/pkg/validate"
)

// Validate checks OpenAPI for path parameter conflicts.
func Validate(doc *openapi3.T) []validate.ValidationError {
	if doc == nil || doc.Paths == nil {
		return nil
	}
	var errs []validate.ValidationError
	for path := range doc.Paths.Map() {
		if hasParamConflict(path) {
			errs = append(errs, validate.ValidationError{
				Rule: "O-1", File: "api/openapi.yaml", Level: "ERROR",
				Message: "path parameter conflict: " + path, SeqIdx: -1,
			})
		}
	}
	return errs
}
