//ff:func feature=crosscheck type=util control=sequence
//ff:what 다이어그램 ID를 DDL 테이블명으로 변환
package crosscheck

import "github.com/jinzhu/inflection"

// diagramIDToTable converts a diagram ID to a DDL table name.
func diagramIDToTable(id string) string {
	return inflection.Plural(id)
}
