//ff:func feature=ssac-gen type=generator control=sequence topic=template-data
//ff:what @call 시퀀스의 패키지명, 함수명, 에러 상태를 templateData에 설정
package generator

import (
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func buildCallModelData(d *templateData, parts []string, seq parser.Sequence, st *validator.SymbolTable) {
	d.PkgName = parts[0]
	if len(parts) > 1 {
		d.FuncMethod = toGoPascal(parts[1])
	}
	d.ErrStatus = resolveCallErrStatus(seq.ErrStatus, st, seq.Model)
}
