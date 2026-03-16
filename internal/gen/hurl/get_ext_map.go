//ff:func feature=gen-hurl type=util control=sequence
//ff:what OpenAPI 확장에서 map을 추출한다
package hurl

import "github.com/getkin/kin-openapi/openapi3"

// getExtMap extracts a map extension value from an OpenAPI operation.
func getExtMap(op *openapi3.Operation, key string) map[string]interface{} {
	if op.Extensions == nil {
		return nil
	}
	v, ok := op.Extensions[key]
	if !ok {
		return nil
	}
	m, ok := v.(map[string]interface{})
	if ok {
		return m
	}
	return nil
}
