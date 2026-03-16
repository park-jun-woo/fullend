//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what SSaC 함수의 첫 번째 @model에서 테이블명 추출
package crosscheck

import (
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// findFirstModelTable extracts the DDL table name from the first @model annotation.
func findFirstModelTable(fn ssacparser.ServiceFunc) string {
	for _, seq := range fn.Sequences {
		if seq.Model == "" {
			continue
		}
		parts := strings.SplitN(seq.Model, ".", 2)
		return modelToTable(parts[0])
	}
	return ""
}
