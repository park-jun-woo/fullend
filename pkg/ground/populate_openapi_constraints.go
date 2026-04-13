//ff:func feature=rule type=loader control=iteration dimension=1
//ff:what populateOpenAPIConstraints — OpenAPI request/response 제약조건을 Ground에 등록
package ground

import (

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateOpenAPIConstraints(g *rule.Ground, fs *fullend.Fullstack) {
	for opID, fields := range fs.RequestConstraints {
		populateConstraintFields(g, "OpenAPI.request."+opID, fields)
	}
	for opID, fields := range fs.ResponseConstraints {
		populateConstraintFields(g, "OpenAPI.response.constraint."+opID, fields)
	}
}
