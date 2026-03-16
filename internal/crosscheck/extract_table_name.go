//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=openapi-ddl
//ff:what CREATE TABLE 문에서 테이블 이름을 추출
package crosscheck

import "strings"

func extractTableName(line string) string {
	parts := strings.Fields(line)
	for i, p := range parts {
		if strings.ToUpper(p) == "TABLE" && i+1 < len(parts) {
			return strings.Trim(parts[i+1], "( ")
		}
	}
	return ""
}
