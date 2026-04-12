//ff:func feature=rule type=rule control=sequence
//ff:what HasPaginateOffset — 시퀀스의 result wrapper가 Page인지 판정
package backend

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func HasPaginateOffset(ctx toulmin.Context, _ toulmin.Specs) (bool, any) {
	claim, _ := ctx.Get("claim")
	seq, _ := claim.(SeqClaim)
	return seq.Seq.Result != nil && seq.Seq.Result.Wrapper == "Page", nil
}
