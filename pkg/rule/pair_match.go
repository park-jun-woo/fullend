//ff:func feature=rule type=rule control=sequence
//ff:what PairMatch — key:value 쌍이 대상에 매칭되는지 검증
package rule

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// PairMatch checks that a (key:value) pair exists in the target.
// claim: string ("key:value" joined). Returns (true, evidence) when NOT found.
func PairMatch(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	s := specs[0].(*PairMatchSpec)
	g, _ := ctx.Get("ground")
	ground := g.(*Ground)
	c, _ := ctx.Get("claim")
	pair, _ := c.(string)
	if pairs, ok := ground.Pairs[s.LookupKey]; ok && pairs[pair] {
		return false, nil
	}
	return true, &Evidence{Rule: s.Rule, Level: s.Level, Ref: pair, Message: s.Message}
}
