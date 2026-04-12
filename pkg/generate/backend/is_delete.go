//ff:func feature=rule type=rule control=sequence
//ff:what IsDelete — 시퀀스 타입이 @delete인지 판정
package backend

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func IsDelete(ctx toulmin.Context, _ toulmin.Specs) (bool, any) {
	claim, _ := ctx.Get("claim")
	seq, _ := claim.(SeqClaim)
	return seq.Type == "delete", nil
}
