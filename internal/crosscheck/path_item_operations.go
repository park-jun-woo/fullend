//ff:func feature=crosscheck type=util control=sequence
//ff:what PathItem에서 모든 Operation 목록 반환
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// pathItemOperations returns all operations from a path item.
func pathItemOperations(pi *openapi3.PathItem) []*openapi3.Operation {
	return []*openapi3.Operation{pi.Get, pi.Post, pi.Put, pi.Delete, pi.Patch}
}
