//ff:func feature=rule type=rule control=sequence
//ff:what HasResult — 시퀀스에 result 변수가 있는지 판정
package backend

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func HasResult(ctx toulmin.Context, _ toulmin.Specs) (bool, any) {
	claim, _ := ctx.Get("claim")
	seq, _ := claim.(SeqClaim)
	return seq.Seq.Result != nil, nil
}
