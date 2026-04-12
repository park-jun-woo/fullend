//ff:func feature=rule type=rule control=sequence
//ff:what IsPut — 시퀀스 타입이 @put인지 판정
package backend

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func IsPut(ctx toulmin.Context, _ toulmin.Specs) (bool, any) {
	claim, _ := ctx.Get("claim")
	seq, _ := claim.(SeqClaim)
	return seq.Type == "put", nil
}
