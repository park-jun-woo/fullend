//ff:func feature=ssac-gen type=generator control=selection topic=template-data
//ff:what 시퀀스 타입별 가드 에러 상태 코드를 templateData에 설정
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func buildErrStatusForGuard(d *templateData, seq ssacparser.Sequence) {
	switch seq.Type {
	case ssacparser.SeqEmpty:
		d.ErrStatus = guardErrStatus(seq.ErrStatus, "http.StatusNotFound")
	case ssacparser.SeqExists:
		d.ErrStatus = guardErrStatus(seq.ErrStatus, "http.StatusConflict")
	case ssacparser.SeqState:
		d.ErrStatus = guardErrStatus(seq.ErrStatus, "http.StatusConflict")
	case ssacparser.SeqAuth:
		d.ErrStatus = guardErrStatus(seq.ErrStatus, "http.StatusForbidden")
	}
}
