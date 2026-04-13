//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkShorthandReverse — OpenAPI field → shorthand @response 커버리지 (X-20)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkShorthandReverse(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil {
		return nil
	}
	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		if fn.Subscribe != nil {
			continue
		}
		if !isShorthandResponse(fn.Sequences) {
			continue
		}
		responseKey := "OpenAPI.response." + fn.Name
		fields := g.Schemas[responseKey]
		if len(fields) <= 1 {
			continue
		}
		errs = append(errs, CrossError{Rule: "X-20", Context: fn.Name, Level: "WARNING",
			Message: "shorthand @response returns single variable but OpenAPI response has multiple fields"})
	}
	return errs
}
