//ff:func feature=ssac-gen type=generator control=sequence topic=response
//ff:what 응답 필드를 templateData에 설정
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func buildResponseData(d *templateData, seq ssacparser.Sequence) {
	d.ResponseFields = seq.Fields
}
