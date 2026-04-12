//ff:func feature=rule type=rule control=sequence
//ff:what HasFKRef — @get의 input이 이전 result 변수의 필드를 참조하는지 판정
package backend

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func HasFKRef(ctx toulmin.Context, _ toulmin.Specs) (bool, any) {
	claim, _ := ctx.Get("claim")
	seq, _ := claim.(SeqClaim)
	return seq.FKRef, nil
}
