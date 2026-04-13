//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what 큐 구독 핸들러 함수를 생성
package ssac

import (
	"fmt"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// generateSubscribeFunc는 큐 구독 핸들러 함수를 생성한다.
func (g *GoTarget) generateSubscribeFunc(sf ssacparser.ServiceFunc, st *validator.SymbolTable) ([]byte, error) {
	pkgName := "service"
	if sf.Feature != "" {
		pkgName = sf.Feature
	}

	bodyBuf := buildSubscribeFuncBody(sf, st, g)

	// message struct 출력: .ssac 파일에 선언된 struct를 함수 본문 앞에 삽입
	structDefs := renderStructDefs(sf.Structs)

	imports := collectSubscribeImports(sf)
	imports = filterUsedImports(imports, bodyBuf.String())

	combined := append(structDefs, bodyBuf.Bytes()...)
	return assembleGoSource(pkgName, imports, combined)
}

// renderStructDefs는 StructInfo 슬라이스를 Go struct 정의로 렌더링한다.
func renderStructDefs(structs []ssacparser.StructInfo) []byte {
	if len(structs) == 0 {
		return nil
	}
	var buf []byte
	for _, s := range structs {
		buf = append(buf, fmt.Sprintf("type %s struct {\n", s.Name)...)
		for _, f := range s.Fields {
			buf = append(buf, fmt.Sprintf("\t%s %s\n", f.Name, f.Type)...)
		}
		buf = append(buf, "}\n\n"...)
	}
	return buf
}
