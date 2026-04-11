//ff:func feature=rule type=rule control=sequence
//ff:what IsSkipped — SSOT kind가 --skip으로 제외되었는지 확인
package rule

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsSkipped checks if a SSOT kind is excluded via --skip flag.
// Ground.Flags["skipped.<kind>"] must be set by the caller.
func IsSkipped(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	s := specs[0].(*SkipSpec)
	g, _ := ctx.Get("ground")
	ground := g.(*Ground)
	return ground.Flags["skipped."+s.Kind], nil
}
