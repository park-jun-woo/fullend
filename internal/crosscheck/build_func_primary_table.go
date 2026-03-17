//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=func-check
//ff:what 함수명을 주 DDL 테이블에 매핑
package crosscheck

import (
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// buildFuncPrimaryTable maps function names to their primary DDL table.
func buildFuncPrimaryTable(funcs []ssacparser.ServiceFunc) map[string]string {
	m := make(map[string]string)
	for _, fn := range funcs {
		if table := findFirstModelTable(fn); table != "" {
			m[fn.Name] = table
		}
	}
	return m
}
