//ff:func feature=rule type=rule control=sequence
//ff:what IsCall — 시퀀스 타입이 @call인지 판정
package backend

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func IsCall(ctx toulmin.Context, _ toulmin.Specs) (bool, any) {
	claim, _ := ctx.Get("claim")
	seq, _ := claim.(SeqClaim)
	return seq.Type == "call", nil
}
