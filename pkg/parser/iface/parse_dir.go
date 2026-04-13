//ff:func feature=iface-parse type=parser control=iteration dimension=1
//ff:what ParseDir — model 디렉토리의 Go 파일을 순회하며 인터페이스 추출
package iface

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/diagnostic"
)

// ParseDir는 주어진 디렉토리의 *.go 파일을 순회하며 Go 인터페이스를 추출한다.
// 디렉토리가 없으면 nil, nil 을 반환한다.
func ParseDir(dir string) ([]Interface, []diagnostic.Diagnostic) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, []diagnostic.Diagnostic{{
			File:    dir,
			Phase:   diagnostic.PhaseParse,
			Level:   diagnostic.LevelError,
			Message: "디렉토리 읽기 실패: " + err.Error(),
		}}
	}

	var ifaces []Interface
	var diags []diagnostic.Diagnostic
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		fileIfaces, fileDiags := ParseFile(path)
		if len(fileDiags) > 0 {
			diags = append(diags, fileDiags...)
			continue
		}
		ifaces = append(ifaces, fileIfaces...)
	}
	return ifaces, diags
}
