//ff:func feature=rule type=rule control=sequence
//ff:what HasSliceResult — 시퀀스의 result 타입이 []T 슬라이스인지 판정
package backend

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func HasSliceResult(ctx toulmin.Context, _ toulmin.Specs) (bool, any) {
	claim, _ := ctx.Get("claim")
	seq, _ := claim.(SeqClaim)
	if seq.Seq.Result == nil {
		return false, nil
	}
	return strings.HasPrefix(seq.Seq.Result.Type, "[]"), nil
}
