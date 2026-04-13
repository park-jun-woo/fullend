//ff:func feature=orchestrator type=rule control=iteration dimension=2
//ff:what DDL NOT NULL 누락 컬럼 감지 — FK DEFAULT 0 센티널 레코드 검사 포함
package orchestrator

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// checkDDLNullableColumns scans DDL files for columns missing NOT NULL.
// PRIMARY KEY columns are implicitly NOT NULL and are excluded.
// Also checks FK + DEFAULT 0 columns for sentinel record (id=0) in referenced table.
func checkDDLNullableColumns(root string, skipSentinel bool) []string {
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
			if msg := checkColumnLine(line, f.tableName, colRe, refRe, tableContents, skipSentinel); msg != "" {
				errs = append(errs, msg)
			}
		}
	}
	return errs
}
