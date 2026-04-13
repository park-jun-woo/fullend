//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=interface-derive
//ff:what @func 호출 결과 변수에 string 타입 단언을 추가한다

package gogin

import "strings"

// addTypeAssertions adds .(string) type assertions for @func results used as string arguments.
func addTypeAssertions(src string, rcv string, funcs []string) string {
	for _, f := range funcs {
		callPattern := rcv + "." + ucFirst(f) + "("
		idx := strings.Index(src, callPattern)
		if idx <= 0 {
			continue
		}
		lineStart := strings.LastIndex(src[:idx], "\n") + 1
		assignLine := strings.TrimSpace(src[lineStart:idx])
		commaIdx := strings.Index(assignLine, ",")
		if commaIdx <= 0 {
			continue
		}
		varName := strings.TrimSpace(assignLine[:commaIdx])
		if varName == "_" || varName == "" {
			continue
		}
		src = strings.ReplaceAll(src, ", "+varName+",", ", "+varName+".(string),")
		src = strings.ReplaceAll(src, ", "+varName+")", ", "+varName+".(string))")
		src = strings.ReplaceAll(src, "("+varName+",", "("+varName+".(string),")
	}
	return src
}
