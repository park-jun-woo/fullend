//ff:func feature=rule type=util control=sequence
//ff:what extractModelFromSeq — 시퀀스의 Model.Method에서 Model 부분 추출
package ssac

import (
	"strings"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func extractModelFromSeq(seq parsessac.Sequence) string {
	if idx := strings.IndexByte(seq.Model, '.'); idx > 0 {
		return seq.Model[:idx]
	}
	return ""
}
