//ff:func feature=crosscheck type=util control=sequence topic=ddl-coverage
//ff:what 단일 SSaC 시퀀스에서 참조되는 DDL 테이블을 수집
package crosscheck

import (
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func collectReferencedTable(seq ssacparser.Sequence, tables map[string]bool) {
	// 패키지 접두사 모델은 DDL 체크 스킵.
	if seq.Package != "" {
		return
	}
	// @model "Course.FindByID" → "courses"
	if seq.Model != "" {
		parts := strings.SplitN(seq.Model, ".", 2)
		if len(parts) >= 1 {
			tables[modelToTable(parts[0])] = true
		}
	}
	// @result type
	if seq.Result != nil && seq.Result.Type != "" {
		typeName := strings.TrimPrefix(seq.Result.Type, "[]")
		typeName = strings.TrimPrefix(typeName, "*")
		if typeName != "" && !primitiveTypes[typeName] {
			tables[modelToTable(typeName)] = true
		}
	}
}
