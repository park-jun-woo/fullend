//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=args-inputs
//ff:what Args 배열을 Go 코드 문자열로 변환
package ssac

import (
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func buildArgsCode(args []ssacparser.Arg) string {
	var parts []string
	for _, a := range args {
		parts = append(parts, argToCode(a))
	}
	return strings.Join(parts, ", ")
}
