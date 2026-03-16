//ff:func feature=crosscheck type=util control=sequence topic=openapi-ddl
//ff:what DDL 테이블명과 타입명의 대응 여부 확인
package crosscheck

import "strings"

// matchTableType checks if a DDL table name corresponds to a type name.
func matchTableType(tableName, typeName string) bool {
	tn := strings.ToLower(typeName)
	tbl := strings.ToLower(tableName)
	return tbl == tn || tbl == tn+"s" || tbl == tn+"es"
}
