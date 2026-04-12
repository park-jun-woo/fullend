//ff:func feature=rule type=rule control=sequence
//ff:what IsGet — 시퀀스 타입이 @get인지 판정
package backend

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsGet activates when sequence type is "get".
func IsGet(ctx toulmin.Context, _ toulmin.Specs) (bool, any) {
	claim, _ := ctx.Get("claim")
	seq, _ := claim.(SeqClaim)
	return seq.Type == "get", nil
}
