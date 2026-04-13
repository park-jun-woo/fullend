//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what compareDDLEnumWithOpenAPI — 단일 DDL CHECK enum 값 목록 ↔ OpenAPI enum 비교
package crosscheck

import (

	"github.com/park-jun-woo/fullend/pkg/fullend"
)

func compareDDLEnumWithOpenAPI(table, col string, ddlVals []string, fs *fullend.Fullstack) []CrossError {
	for _, fields := range fs.ResponseConstraints {
		fc, ok := fields[col]
		if !ok || len(fc.Enum) == 0 {
			continue
		}
		return matchDDLValsToEnum(table, col, ddlVals, fc.Enum)
	}
	return nil
}
