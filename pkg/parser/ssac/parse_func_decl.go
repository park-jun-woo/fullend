//ff:func feature=ssac-parse type=parser control=sequence
//ff:what AST 함수 선언에서 ServiceFunc를 추출
package parser

import (
	"fmt"
	"go/ast"
	"path/filepath"
)

// parseFuncDecl은 AST 함수 선언에서 ServiceFunc를 추출한다.
func parseFuncDecl(fn *ast.FuncDecl, f *ast.File, path string, imports []string, structs []StructInfo) (*ServiceFunc, error) {
	comments := collectFuncComments(f, fn.Pos())

	sequences, err := parseComments(comments)
	if err != nil {
		return nil, fmt.Errorf("%s:%s — %w", filepath.Base(path), fn.Name.Name, err)
	}
	if len(sequences) == 0 {
		return nil, nil
	}

	sf := ServiceFunc{
		Name:     fn.Name.Name,
		FileName: filepath.Base(path),
		Imports:  imports,
		Structs:  structs,
		Param:    extractParamInfo(fn),
	}

	// @subscribe 추출: 시퀀스가 아닌 함수 메타데이터
	sf.Sequences = filterSubscribe(&sf, sequences)

	return &sf, nil
}
