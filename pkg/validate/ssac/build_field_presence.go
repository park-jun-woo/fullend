//ff:func feature=rule type=util control=sequence
//ff:what buildFieldPresence — 시퀀스에서 필드 존재 여부 맵 구성
package ssac

import parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func buildFieldPresence(seq parsessac.Sequence) map[string]bool {
	return map[string]bool{
		"Model":      seq.Model != "",
		"Result":     seq.Result != nil,
		"Args":       len(seq.Args) > 0,
		"Target":     seq.Target != "",
		"Message":    seq.Message != "",
		"DiagramID":  seq.DiagramID != "",
		"Inputs":     len(seq.Inputs) > 0,
		"Transition": seq.Transition != "",
		"Action":     seq.Action != "",
		"Resource":   seq.Resource != "",
		"Topic":      seq.Topic != "",
		"Payload":    len(seq.Fields) > 0,
	}
}
