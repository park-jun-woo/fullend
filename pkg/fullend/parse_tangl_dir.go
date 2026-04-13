//ff:func feature=orchestrator type=parser control=iteration dimension=1
//ff:what TANGL .md 파일들을 파싱·검증하여 유효한 File 목록 반환
package fullend

import (
	"os"
	"path/filepath"
	"strings"

	tanglparser "github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

func parseTanglDir(dir string) []*tanglparser.File {
	entries, _ := os.ReadDir(dir)
	var files []*tanglparser.File
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		f := parseTanglFile(filepath.Join(dir, e.Name()))
		if f != nil {
			files = append(files, f)
		}
	}
	return files
}
