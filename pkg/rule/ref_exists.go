//ff:func feature=rule type=rule control=sequence
//ff:what RefExists — 참조 이름이 대상에 존재하는지 검증
package rule

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// RefExists checks that a referenced name exists in the target.
// claim: string (reference name). Returns (true, evidence) when NOT found.
func RefExists(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	s := specs[0].(*RefExistsSpec)
	g, _ := ctx.Get("ground")
	ground := g.(*Ground)
	c, _ := ctx.Get("claim")
	name, _ := c.(string)
	if names, ok := ground.Lookup[s.LookupKey]; ok && names[name] {
		return false, nil
	}
	return true, &Evidence{Rule: s.Rule, Level: s.Level, Ref: name, Message: s.Message}
}
