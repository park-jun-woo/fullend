//ff:type feature=gen-hurl type=model
//ff:what scenarioStep represents a single HTTP request in the Hurl scenario.
package hurl

import "github.com/getkin/kin-openapi/openapi3"

// scenarioStep represents a single HTTP request in the Hurl scenario.
type scenarioStep struct {
	OperationID string
	Method      string
	Path        string
	Operation   *openapi3.Operation
	PathDepth   int  // number of path segments (for ordering)
	IsAuth      bool // register/login
}
