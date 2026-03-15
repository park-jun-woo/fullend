//ff:func feature=symbol type=util
//ff:what SQL 파일명을 모델명으로 변환한다 (reservations.sql → Reservation)
package validator

import (
	"strings"

	"github.com/ettle/strcase"
	"github.com/jinzhu/inflection"
)

// sqlFileToModel은 "reservations.sql" → "Reservation" 변환한다.
func sqlFileToModel(filename string) string {
	name := strings.TrimSuffix(filename, ".sql")
	singular := inflection.Singular(name)
	return strcase.ToGoPascal(singular)
}
