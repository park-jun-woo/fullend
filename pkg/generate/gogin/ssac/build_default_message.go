//ff:func feature=ssac-gen type=generator control=sequence topic=template-data
//ff:what 시퀀스의 메시지가 비어있으면 기본 메시지를 설정
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func buildDefaultMessage(d *templateData, seq ssacparser.Sequence) {
	if d.Message == "" {
		d.Message = defaultMessage(seq)
	}
}
