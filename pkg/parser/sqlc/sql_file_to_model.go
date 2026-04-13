//ff:func feature=sqlc-parse type=util control=sequence
//ff:what sqlFileToModel — "reservations.sql" → "Reservation" (단수화 + PascalCase)
package sqlc

import (
	"strings"

	"github.com/ettle/strcase"
	"github.com/jinzhu/inflection"
)

func sqlFileToModel(filename string) string {
	name := strings.TrimSuffix(filename, ".sql")
	singular := inflection.Singular(name)
	return strcase.ToGoPascal(singular)
}
