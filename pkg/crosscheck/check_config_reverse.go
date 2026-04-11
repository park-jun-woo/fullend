//ff:func feature=crosscheck type=rule control=sequence
//ff:what checkConfigReverse — Config middleware → OpenAPI security 역방향 (X-51), claims 커버리지 (X-54)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkConfigReverse(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.Manifest == nil {
		return nil
	}
	var errs []CrossError

	// X-51: OpenAPI security → Config middleware
	if fs.OpenAPIDoc != nil {
		graph51 := toulmin.NewGraph("security-middleware")
		graph51.Rule(rule.PairMatch).With(&rule.PairMatchSpec{
			BaseSpec:  rule.BaseSpec{Rule: "X-51", Level: "ERROR", Message: "OpenAPI securityScheme not in Config middleware"},
			LookupKey: "Config.middleware",
		})
		for name := range g.Lookup["OpenAPI.security"] {
			errs = append(errs, evalRef(graph51, g, name, name)...)
		}
	}

	// X-54: Config claims → Rego usage coverage
	if len(fs.ParsedPolicies) > 0 {
		graph54 := toulmin.NewGraph("claims-coverage")
		graph54.Rule(rule.CoverageCheck).With(&rule.CoverageCheckSpec{
			BaseSpec:  rule.BaseSpec{Rule: "X-54", Level: "WARNING", Message: "Config claim not referenced in Rego"},
			LookupKey: "Rego.claims",
		})
		for field := range g.Lookup["Config.claims"] {
			errs = append(errs, evalRef(graph54, g, field, field)...)
		}
	}

	return errs
}
