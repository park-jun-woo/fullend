//ff:func feature=ssac-gen type=generator control=sequence topic=publish
//ff:what publish 시퀀스의 토픽, 페이로드, 옵션을 templateData에 설정
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func buildPublishData(d *templateData, seq ssacparser.Sequence) {
	if seq.Type != ssacparser.SeqPublish {
		return
	}
	d.Topic = seq.Topic
	d.InputFields = buildPublishPayload(seq.Inputs)
	d.OptionCode = buildPublishOptions(seq.Options)
}
