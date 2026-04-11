//ff:func feature=rule type=rule control=sequence
//ff:what IsDTO — 현재 타입이 @dto 마크인지 확인
package rule

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsDTO checks if the current type is marked // @dto.
// Ground.Flags["dto"] must be set by the caller.
func IsDTO(ctx toulmin.Context, _ toulmin.Specs) (bool, any) {
	g, _ := ctx.Get("ground")
	ground := g.(*Ground)
	return ground.Flags["dto"], nil
}
