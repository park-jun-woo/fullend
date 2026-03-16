//ff:func feature=ssac-gen type=generator control=sequence
//ff:what publish 시퀀스의 토픽, 페이로드, 옵션을 templateData에 설정
package generator

import "github.com/geul-org/fullend/internal/ssac/parser"

func buildPublishData(d *templateData, seq parser.Sequence) {
	if seq.Type != parser.SeqPublish {
		return
	}
	d.Topic = seq.Topic
	d.InputFields = buildPublishPayload(seq.Inputs)
	d.OptionCode = buildPublishOptions(seq.Options)
}
