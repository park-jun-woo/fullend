//ff:func feature=rule type=loader control=iteration dimension=2
//ff:what populateRequestSchemas — OpenAPI requestBody 필드 제약을 g.ReqSchemas 로 변환
package ground

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	oapiparser "github.com/park-jun-woo/fullend/pkg/parser/openapi"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateRequestSchemas(g *rule.Ground, fs *fullend.Fullstack) {
	if g.ReqSchemas == nil {
		g.ReqSchemas = make(map[string]rule.RequestSchemaInfo)
	}
	for opID, fields := range fs.RequestConstraints {
		dst := make(map[string]rule.FieldConstraint, len(fields))
		for name, fc := range fields {
			dst[name] = toRuleFieldConstraint(fc)
		}
		g.ReqSchemas[opID] = rule.RequestSchemaInfo{Fields: dst}
	}
}

func toRuleFieldConstraint(src oapiparser.FieldConstraint) rule.FieldConstraint {
	return rule.FieldConstraint{
		Required:  src.Required,
		Format:    src.Format,
		MinLength: src.MinLength,
		MaxLength: src.MaxLength,
		Enum:      src.Enum,
	}
}
