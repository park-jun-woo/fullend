//ff:func feature=crosscheck type=util control=sequence
//ff:what 모델 이름을 DDL 테이블 이름으로 변환
package crosscheck

import "github.com/jinzhu/inflection"

// primitiveTypes are Go types that never map to DDL tables.
var primitiveTypes = map[string]bool{
	"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"float32": true, "float64": true,
	"string": true, "bool": true, "byte": true, "rune": true,
	"error": true, "any": true,
}

// modelToTable converts a model name to a table name.
// e.g. "User" → "users", "Reservation" → "reservations", "Room" → "rooms"
func modelToTable(model string) string {
	return inflection.Plural(pascalToSnake(model))
}
