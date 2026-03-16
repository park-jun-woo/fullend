//ff:func feature=symbol type=util control=iteration dimension=1 topic=ddl
//ff:what 컬럼 정의에서 인라인 REFERENCES를 파싱한다
package validator

import "strings"

// parseInlineFK는 컬럼 정의에서 인라인 REFERENCES를 파싱한다.
// e.g. "user_id BIGINT NOT NULL REFERENCES users(id)"
func parseInlineFK(colName string, parts []string) (ForeignKey, bool) {
	for i, p := range parts {
		if strings.ToUpper(p) != "REFERENCES" || i+1 >= len(parts) {
			continue
		}
		ref := parts[i+1]
		ref = strings.TrimSuffix(ref, ",")
		refTable, refCol := parseRef(ref)
		if refTable != "" {
			return ForeignKey{Column: colName, RefTable: refTable, RefColumn: refCol}, true
		}
	}
	return ForeignKey{}, false
}
