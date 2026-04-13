//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkJWTClaims — JWT @call input → Config claims 필드 검증 (X-73)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkJWTClaims(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.Manifest == nil || fs.Manifest.Backend.Auth == nil {
		return nil
	}
	graph := toulmin.NewGraph("jwt-claims")
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-73", Level: "ERROR", Message: "JWT @call input not in Config claims"},
		LookupKey: "Config.claims",
	})

	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		errs = append(errs, checkJWTClaimSeqs(graph, g, fn.Name, fn.Sequences)...)
	}
	return errs
}
