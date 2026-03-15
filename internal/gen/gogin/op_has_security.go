//ff:func feature=gen-gogin type=util
//ff:what returns true if an OpenAPI operation has a non-empty security requirement

package gogin

import "github.com/getkin/kin-openapi/openapi3"

// opHasSecurity returns true if an OpenAPI operation has a non-empty security requirement.
func opHasSecurity(op *openapi3.Operation) bool {
	if op.Security == nil {
		return false
	}
	// security: [] means explicitly no auth.
	// security: [{bearerAuth: []}] means auth required.
	return len(*op.Security) > 0
}
