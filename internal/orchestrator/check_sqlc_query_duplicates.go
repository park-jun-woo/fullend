//ff:func feature=orchestrator type=rule control=iteration dimension=2
//ff:what sqlc 쿼리 이름 중복 감지 — db/queries/*.sql 스캔
package orchestrator

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// checkSqlcQueryDuplicates scans db/queries/*.sql for duplicate -- name: entries.
func checkSqlcQueryDuplicates(root string) []string {
	queriesDir := filepath.Join(root, "db", "queries")
	entries, err := os.ReadDir(queriesDir)
	if err != nil {
		return nil
	}

	nameRe := regexp.MustCompile(`^--\s*name:\s*(\w+)\s+:(\w+)`)
	// nameToFiles maps query name -> list of filenames where it appears.
	nameToFiles := make(map[string][]string)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		f, err := os.Open(filepath.Join(queriesDir, entry.Name()))
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if m := nameRe.FindStringSubmatch(scanner.Text()); m != nil {
				nameToFiles[m[1]] = append(nameToFiles[m[1]], entry.Name())
			}
		}
		f.Close()
	}

	var errs []string
	for name, files := range nameToFiles {
		if len(files) > 1 {
			errs = append(errs, fmt.Sprintf(
				"db/queries: %q 이름이 중복됩니다 (%s) — sqlc는 전역 네임스페이스이므로 ModelPrefix를 붙이세요 (예: User%s, Gig%s)",
				name, strings.Join(files, ", "), name, name))
		}
	}
	return errs
}
