//ff:func feature=ssac-gen type=generator control=sequence topic=template-data
//ff:what 상태(state)와 인가(auth) 관련 필드를 templateData에 설정
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func buildStateAuth(d *templateData, seq ssacparser.Sequence) {
	d.DiagramID = seq.DiagramID
	d.Transition = seq.Transition
	d.Action = seq.Action
	d.Resource = seq.Resource
}
