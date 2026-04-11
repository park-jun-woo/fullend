//ff:func feature=rule type=rule control=sequence
//ff:what ForbiddenRef — 이름이 금지 목록에 없는지 검증
package rule

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ForbiddenRef checks that a name is NOT in the forbidden set.
// claim: string (name). Returns (true, evidence) when name IS found (violation).
func ForbiddenRef(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	s := specs[0].(*ForbiddenRefSpec)
	g, _ := ctx.Get("ground")
	ground := g.(*Ground)
	c, _ := ctx.Get("claim")
	name, _ := c.(string)
	if names, ok := ground.Lookup[s.LookupKey]; ok && names[name] {
		return true, &Evidence{Rule: s.Rule, Level: s.Level, Ref: name, Message: s.Message}
	}
	return false, nil
}
