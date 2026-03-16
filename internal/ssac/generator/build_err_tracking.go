//ff:func feature=ssac-gen type=generator control=selection
//ff:what 시퀀스 타입별 err 변수 선언 추적 및 FirstErr 플래그 설정
package generator

import "github.com/geul-org/fullend/internal/ssac/parser"

func buildErrTracking(d *templateData, seq parser.Sequence, errDeclared *bool) {
	d.ErrDeclared = *errDeclared

	switch seq.Type {
	case parser.SeqGet, parser.SeqPost:
		d.FirstErr = true
		*errDeclared = true
	case parser.SeqAuth:
		if !*errDeclared {
			d.FirstErr = true
			*errDeclared = true
		}
	case parser.SeqCall:
		if seq.Result != nil || !*errDeclared {
			d.FirstErr = true
			*errDeclared = true
		}
	case parser.SeqPut, parser.SeqDelete, parser.SeqPublish:
		if !*errDeclared {
			d.FirstErr = true
			*errDeclared = true
		}
	}
}
