//ff:func feature=orchestrator type=rule
//ff:what DDL NOT NULL 누락 컬럼 감지 — FK DEFAULT 0 센티널 레코드 검사 포함
package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// checkDDLNullableColumns scans DDL files for columns missing NOT NULL.
// PRIMARY KEY columns are implicitly NOT NULL and are excluded.
// Also checks FK + DEFAULT 0 columns for sentinel record (id=0) in referenced table.
func checkDDLNullableColumns(root string) []string {
	dbDir := filepath.Join(root, "db")
	entries, err := os.ReadDir(dbDir)
	if err != nil {
		return nil
	}

	createRe := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(\w+)`)
	colRe := regexp.MustCompile(`^(\w+)\s+\w+`)
	refRe := regexp.MustCompile(`(?i)REFERENCES\s+(\w+)`)

	// 1단계: 모든 DDL 파일 내용을 테이블별로 수집.
	tableContents := make(map[string]string) // tableName → 파일 전체 내용
	type fileInfo struct {
		tableName string
		content   string
	}
	var files []fileInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dbDir, entry.Name()))
		if err != nil {
			continue
		}
		content := string(data)
		tableMatch := createRe.FindStringSubmatch(content)
		if tableMatch == nil {
			continue
		}
		tableName := tableMatch[1]
		tableContents[tableName] = content
		files = append(files, fileInfo{tableName: tableName, content: content})
	}

	// 2단계: 컬럼별 NOT NULL 체크 + FK DEFAULT 0 센티널 체크.
	var errs []string
	for _, f := range files {
		for _, line := range strings.Split(f.content, "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "--") || strings.HasPrefix(strings.ToUpper(trimmed), "CREATE") || strings.HasPrefix(trimmed, ")") {
				continue
			}
			upper := strings.ToUpper(trimmed)
			// Skip non-DDL statements (INSERT, ON, etc.).
			if strings.HasPrefix(upper, "INSERT") || strings.HasPrefix(upper, "ON ") || strings.HasPrefix(upper, "VALUES") {
				continue
			}
			if strings.HasPrefix(upper, "PRIMARY KEY") || strings.HasPrefix(upper, "UNIQUE") || strings.HasPrefix(upper, "CHECK") || strings.HasPrefix(upper, "FOREIGN KEY") || strings.HasPrefix(upper, "CONSTRAINT") {
				continue
			}
			m := colRe.FindStringSubmatch(trimmed)
			if m == nil {
				continue
			}
			colName := m[1]
			if strings.Contains(upper, "PRIMARY KEY") || strings.Contains(upper, "NOT NULL") {
				// FK + DEFAULT 0 패턴: 참조 대상 테이블에 id=0 센티널 레코드 확인.
				if strings.Contains(upper, "DEFAULT 0") && strings.Contains(upper, "REFERENCES") {
					refMatch := refRe.FindStringSubmatch(trimmed)
					if refMatch != nil {
						refTable := refMatch[1]
						if refContent, ok := tableContents[refTable]; ok {
							if !hasSentinelInsert(refContent, refTable) {
								errs = append(errs, fmt.Sprintf("DDL: 테이블 %q 컬럼 %q — FK + DEFAULT 0이지만 참조 대상 %q에 id=0 센티널 레코드가 없습니다. INSERT INTO %s (id, ...) VALUES (0, ...) ON CONFLICT DO NOTHING; 을 추가하세요", f.tableName, colName, refTable, refTable))
							}
						}
					}
				}
				continue
			}
			errs = append(errs, fmt.Sprintf("DDL: 테이블 %q 컬럼 %q — NOT NULL이 없습니다. NOT NULL DEFAULT 값을 지정하세요", f.tableName, colName))
		}
	}
	return errs
}
