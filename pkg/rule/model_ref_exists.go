//ff:func feature=rule type=rule control=sequence
//ff:what ModelRefExists — 참조된 모델이 Ground.Models 에 존재하는지 검증
package rule

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ModelRefExists checks that a referenced model name exists in Ground.Models.
// claim: string (model name). Returns (true, evidence) when NOT found.
// 구조적 g.Models 기반이라 legacy Lookup["SymbolTable.model"] 을 대체한다.
func ModelRefExists(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	s := specs[0].(*ModelRefExistsSpec)
	g, _ := ctx.Get("ground")
	ground := g.(*Ground)
	c, _ := ctx.Get("claim")
	name, _ := c.(string)
	if _, ok := ground.Models[name]; ok {
		return false, nil
	}
	return true, &Evidence{Rule: s.Rule, Level: s.Level, Ref: name, Message: s.Message}
}
