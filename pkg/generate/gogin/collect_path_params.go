//ff:func feature=gen-gogin type=util control=iteration dimension=2
//ff:what collects path parameter info from OpenAPI path item and operation parameters

package gogin

import "github.com/getkin/kin-openapi/openapi3"

// collectPathParams collects path parameter info from a path item and operation.
func collectPathParams(pathItem *openapi3.PathItem, op *openapi3.Operation) []pathParamInfo {
	var params []pathParamInfo
	for _, src := range [2]openapi3.Parameters{pathItem.Parameters, op.Parameters} {
		for _, p := range src {
			if p.Value == nil || p.Value.In != "path" {
				continue
			}
			params = append(params, pathParamInfo{
				Name:   p.Value.Name,
				GoName: snakeToGo(p.Value.Name),
				IsInt:  p.Value.Schema != nil && p.Value.Schema.Value != nil && p.Value.Schema.Value.Type != nil && ((*p.Value.Schema.Value.Type)[0] == "integer"),
			})
		}
	}
	return params
}
