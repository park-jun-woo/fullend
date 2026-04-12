//ff:func feature=rule type=rule control=sequence
//ff:what HasPaginateCursor — 시퀀스의 result wrapper가 Cursor인지 판정
package backend

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func HasPaginateCursor(ctx toulmin.Context, _ toulmin.Specs) (bool, any) {
	claim, _ := ctx.Get("claim")
	seq, _ := claim.(SeqClaim)
	return seq.Seq.Result != nil && seq.Seq.Result.Wrapper == "Cursor", nil
}
