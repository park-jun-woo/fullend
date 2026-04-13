//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkOpenAPIRequired — SSaC used fields → OpenAPI required 포함 여부 (X-66)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkOpenAPIRequired(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil {
		return nil
	}
	var errs []CrossError
	for opID, fields := range fs.RequestConstraints {
		required := g.Schemas["OpenAPI.request."+opID+".required"]
		if len(required) == 0 {
			continue
		}
		errs = append(errs, evalRequiredCoverage(g, opID, required, fields)...)
	}
	return errs
}
