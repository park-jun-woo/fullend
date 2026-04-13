//ff:func feature=orchestrator type=parser control=sequence
//ff:what 단일 TANGL .md 파일을 파싱·검증하여 유효하면 File 반환
package fullend

import (
	tanglparser "github.com/park-jun-woo/toulmin/pkg/tangl/parser"
	tanglvalidate "github.com/park-jun-woo/toulmin/pkg/tangl/validate"
)

func parseTanglFile(path string) *tanglparser.File {
	f, err := tanglparser.Parse(path)
	if err != nil {
		return nil
	}
	if err := tanglvalidate.Validate(f); err != nil {
		return nil
	}
	return f
}
