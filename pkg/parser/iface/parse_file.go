//ff:func feature=iface-parse type=parser control=iteration dimension=1
//ff:what ParseFile — 단일 Go 파일에서 인터페이스 선언을 추출
package iface

import (
	"go/parser"
	"go/token"

	"github.com/park-jun-woo/fullend/pkg/diagnostic"
)

// ParseFile는 단일 Go 소스 파일을 파싱해 인터페이스 선언을 추출한다.
func ParseFile(path string) ([]Interface, []diagnostic.Diagnostic) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, []diagnostic.Diagnostic{{
			File:    path,
			Phase:   diagnostic.PhaseParse,
			Level:   diagnostic.LevelError,
			Message: "Go 파일 파싱 실패: " + err.Error(),
		}}
	}
	return extractInterfaces(f), nil
}
