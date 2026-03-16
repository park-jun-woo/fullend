//ff:func feature=symbol type=loader control=iteration dimension=1
//ff:what db/ 디렉토리의 DDL .sql 파일에서 CREATE TABLE 문의 컬럼 타입을 추출한다
package validator

import (
	"os"
	"path/filepath"
	"strings"
)

// loadDDL은 db/ 디렉토리의 DDL .sql 파일에서 CREATE TABLE 문의 컬럼 타입을 추출한다.
func (st *SymbolTable) loadDDL(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return err
		}

		parseDDLTables(string(data), st.DDLTables)
	}
	return nil
}
