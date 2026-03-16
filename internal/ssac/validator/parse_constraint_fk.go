//ff:func feature=symbol type=util control=sequence
//ff:what 독립 FOREIGN KEY 절을 파싱한다
package validator

import "strings"

// parseConstraintFK는 독립 FOREIGN KEY 절을 파싱한다.
// e.g. "CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)"
// e.g. "FOREIGN KEY (user_id) REFERENCES users(id)"
func parseConstraintFK(line string) (ForeignKey, bool) {
	upper := strings.ToUpper(line)
	fkIdx := strings.Index(upper, "FOREIGN KEY")
	refIdx := strings.Index(upper, "REFERENCES")
	if fkIdx < 0 || refIdx < 0 {
		return ForeignKey{}, false
	}

	// FOREIGN KEY (col) 부분에서 컬럼 추출
	between := line[fkIdx+len("FOREIGN KEY") : refIdx]
	col := extractParenContent(between)
	if col == "" {
		return ForeignKey{}, false
	}

	// REFERENCES table(col) 부분
	after := strings.TrimSpace(line[refIdx+len("REFERENCES"):])
	after = strings.TrimSuffix(after, ",")
	refTable, refCol := parseRef(after)
	if refTable == "" {
		return ForeignKey{}, false
	}

	return ForeignKey{Column: col, RefTable: refTable, RefColumn: refCol}, true
}
