//ff:func feature=contract type=util control=iteration dimension=1
//ff:what 단일 시퀀스의 Args에서 request 소스 필드명을 수집한다
package contract

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

// collectReqFieldsFromSeq returns request field names from a single sequence's args.
func collectReqFieldsFromSeq(args []ssacparser.Arg) []string {
	var fields []string
	for _, arg := range args {
		if arg.Source == "request" {
			fields = append(fields, arg.Field)
		}
	}
	return fields
}
