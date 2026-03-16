//ff:func feature=ssac-gen type=generator control=selection
//ff:what 시퀀스 타입별 가드 에러 상태 코드를 templateData에 설정
package generator

import "github.com/geul-org/fullend/internal/ssac/parser"

func buildErrStatusForGuard(d *templateData, seq parser.Sequence) {
	switch seq.Type {
	case parser.SeqEmpty:
		d.ErrStatus = guardErrStatus(seq.ErrStatus, "http.StatusNotFound")
	case parser.SeqExists:
		d.ErrStatus = guardErrStatus(seq.ErrStatus, "http.StatusConflict")
	case parser.SeqState:
		d.ErrStatus = guardErrStatus(seq.ErrStatus, "http.StatusConflict")
	case parser.SeqAuth:
		d.ErrStatus = guardErrStatus(seq.ErrStatus, "http.StatusForbidden")
	}
}
