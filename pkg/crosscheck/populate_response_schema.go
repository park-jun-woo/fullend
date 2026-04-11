//ff:func feature=crosscheck type=loader control=iteration dimension=1
//ff:what populateResponseSchema — operationId별 OpenAPI response 필드를 Ground.Schemas에 등록
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateResponseSchema(g *rule.Ground, opID string, op *openapi3.Operation) {
	if op.Responses == nil {
		return
	}
	for code, resp := range op.Responses.Map() {
		if len(code) == 0 || code[0] != '2' {
			continue
		}
		if resp.Value == nil || resp.Value.Content == nil {
			continue
		}
		ct := resp.Value.Content.Get("application/json")
		if ct == nil || ct.Schema == nil || ct.Schema.Value == nil {
			continue
		}
		var fields []string
		for name := range ct.Schema.Value.Properties {
			fields = append(fields, name)
		}
		g.Schemas["OpenAPI.response."+opID] = fields
		return
	}
}
