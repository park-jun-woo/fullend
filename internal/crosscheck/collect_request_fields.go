//ff:func feature=crosscheck type=util control=iteration dimension=2 topic=openapi-ddl
//ff:what ServiceFunc의 시퀀스에서 request.X 참조 필드를 수집
package crosscheck

import (
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func collectRequestFields(fn ssacparser.ServiceFunc) map[string]bool {
	fields := make(map[string]bool)
	for _, seq := range fn.Sequences {
		for _, arg := range seq.Args {
			if arg.Source == "request" && arg.Field != "" {
				fields[arg.Field] = true
			}
		}
		for _, val := range seq.Inputs {
			if strings.HasPrefix(val, "request.") {
				fields[strings.TrimPrefix(val, "request.")] = true
			}
		}
		for _, val := range seq.Fields {
			if strings.HasPrefix(val, "request.") {
				fields[strings.TrimPrefix(val, "request.")] = true
			}
		}
	}
	return fields
}
