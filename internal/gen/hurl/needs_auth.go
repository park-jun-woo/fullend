//ff:func feature=gen-hurl type=util
//ff:what 오퍼레이션이 인증을 필요로 하는지 확인한다
package hurl

import "github.com/getkin/kin-openapi/openapi3"

// needsAuth returns true if the operation requires authentication.
func needsAuth(op *openapi3.Operation) bool {
	return op.Security != nil && len(*op.Security) > 0
}
