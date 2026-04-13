//ff:func feature=ssac-gen type=generator control=selection topic=template-data
//ff:what 시퀀스 타입별 err 변수 선언 추적 및 FirstErr 플래그 설정
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func buildErrTracking(d *templateData, seq ssacparser.Sequence, errDeclared *bool) {
	d.ErrDeclared = *errDeclared

	switch seq.Type {
	case ssacparser.SeqGet, ssacparser.SeqPost:
		d.FirstErr = true
		*errDeclared = true
	case ssacparser.SeqAuth:
		if !*errDeclared {
			d.FirstErr = true
			*errDeclared = true
		}
	case ssacparser.SeqCall:
		if seq.Result != nil || !*errDeclared {
			d.FirstErr = true
			*errDeclared = true
		}
	case ssacparser.SeqPut, ssacparser.SeqDelete, ssacparser.SeqPublish:
		if !*errDeclared {
			d.FirstErr = true
			*errDeclared = true
		}
	}
}
