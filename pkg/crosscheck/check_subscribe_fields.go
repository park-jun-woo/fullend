//ff:func feature=crosscheck type=rule control=sequence
//ff:what checkSubscribeFields — @subscribe message 필드 → @publish payload 필드 매칭 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkSubscribeFields(fn ssac.ServiceFunc, pubPayloads map[string][]string, g *rule.Ground) []CrossError {
	pubFields, ok := pubPayloads[fn.Subscribe.Topic]
	if !ok {
		return nil
	}
	subFields := collectMessageFields(fn.Structs, fn.Subscribe.MessageType)
	if len(subFields) == 0 {
		return nil
	}
	localG := shallowCopyGround(g)
	localG.Schemas["_pub"] = pubFields
	return evalSchemaMatch(localG, subFields, pubFields, "X-59", fn.Name)
}
