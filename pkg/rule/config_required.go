//ff:func feature=rule type=rule control=sequence
//ff:what ConfigRequired — 설정 키가 존재하는지 검증
package rule

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ConfigRequired checks that a config key is set.
// claim: ignored. Returns (true, evidence) when the key is NOT set.
func ConfigRequired(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	s := specs[0].(*ConfigRequiredSpec)
	g, _ := ctx.Get("ground")
	ground := g.(*Ground)
	if ground.Config[s.ConfigKey] {
		return false, nil
	}
	return true, &Evidence{Rule: s.Rule, Level: s.Level, Message: s.Message}
}
