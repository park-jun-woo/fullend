//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkColMaxLength — 단일 컬럼의 VARCHAR(n) ↔ OpenAPI maxLength 비교
package crosscheck

import (
	"strconv"

	"github.com/park-jun-woo/fullend/pkg/fullend"
)

func checkColMaxLength(table, col string, vLen int, fs *fullend.Fullstack) []CrossError {
	for _, fields := range fs.ResponseConstraints {
		fc, ok := fields[col]
		if !ok {
			continue
		}
		if fc.MaxLength == nil {
			return []CrossError{{Rule: "X-67", Context: table + "." + col, Level: "WARNING",
				Message: "DDL VARCHAR(" + strconv.Itoa(vLen) + ") but OpenAPI has no maxLength"}}
		}
		if *fc.MaxLength > vLen {
			return []CrossError{{Rule: "X-70", Context: table + "." + col, Level: "WARNING",
				Message: "OpenAPI maxLength " + strconv.Itoa(*fc.MaxLength) + " exceeds DDL VARCHAR(" + strconv.Itoa(vLen) + ")"}}
		}
		return nil
	}
	return nil
}
