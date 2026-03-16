//ff:func feature=ssac-gen type=generator control=sequence topic=response
//ff:what 응답 필드를 templateData에 설정
package generator

import "github.com/geul-org/fullend/internal/ssac/parser"

func buildResponseData(d *templateData, seq parser.Sequence) {
	d.ResponseFields = seq.Fields
}
