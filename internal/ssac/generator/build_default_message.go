//ff:func feature=ssac-gen type=generator control=sequence
//ff:what 시퀀스의 메시지가 비어있으면 기본 메시지를 설정
package generator

import "github.com/geul-org/fullend/internal/ssac/parser"

func buildDefaultMessage(d *templateData, seq parser.Sequence) {
	if d.Message == "" {
		d.Message = defaultMessage(seq)
	}
}
