//ff:func feature=ssac-gen type=generator control=sequence topic=template-data
//ff:what query + 리스트 반환 시 HasTotal 플래그를 설정
package ssac

import (
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func buildHasTotal(d *templateData, seq ssacparser.Sequence) {
	if hasQueryInput(seq.Inputs) && seq.Result != nil && strings.HasPrefix(seq.Result.Type, "[]") && seq.Result.Wrapper == "" {
		d.HasTotal = true
	}
}
