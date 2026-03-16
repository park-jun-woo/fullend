//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=sensitive
//ff:what SQL 텍스트에서 @sensitive/@nosensitive 컬럼 태그를 추출
package crosscheck

import "strings"

func parseSensitiveSQL(content string, sensitive, nosensitive map[string]map[string]bool) {
	lines := strings.Split(content, "\n")
	var currentTable string

	for _, line := range lines {
		currentTable = parseSensitiveLine(line, currentTable, sensitive, nosensitive)
	}
}
