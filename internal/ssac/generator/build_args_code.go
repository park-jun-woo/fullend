//ff:func feature=ssac-gen type=generator control=iteration dimension=1
//ff:what Args 배열을 Go 코드 문자열로 변환
package generator

import (
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

func buildArgsCode(args []parser.Arg) string {
	var parts []string
	for _, a := range args {
		parts = append(parts, argToCode(a))
	}
	return strings.Join(parts, ", ")
}
