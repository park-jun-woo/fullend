//ff:func feature=rule type=rule control=sequence
//ff:what TypeMatch — 소스 타입이 대상 타입과 일치하는지 검증
package rule

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// TypeMatch checks that the source type matches the target type.
// claim: *TypeClaim. Returns (true, evidence) when types do NOT match.
func TypeMatch(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	s := specs[0].(*TypeMatchSpec)
	g, _ := ctx.Get("ground")
	ground := g.(*Ground)
	c, _ := ctx.Get("claim")
	tc, _ := c.(*TypeClaim)
	targetType, ok := ground.Types[s.LookupKey+"."+tc.Name]
	if !ok {
		return false, nil
	}
	if tc.SourceType == targetType {
		return false, nil
	}
	return true, &Evidence{Rule: s.Rule, Level: s.Level, Ref: tc.Name, Message: s.Message}
}
