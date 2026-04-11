//ff:func feature=crosscheck type=loader control=selection
//ff:what populateSSaCSeq — 개별 시퀀스에서 auth/call/model/publish 정보 추출
package crosscheck

import (
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateSSaCSeq(g *rule.Ground, funcName string, seq ssac.Sequence,
	authPairs, callRefs, modelRefs, pubTopics rule.StringSet) {

	switch seq.Type {
	case "auth":
		authPairs[seq.Action+":"+seq.Resource] = true
	case "call":
		if idx := strings.IndexByte(seq.Model, '.'); idx > 0 {
			callRefs[strings.ToLower(seq.Model)] = true
		}
	case "publish":
		pubTopics[seq.Topic] = true
	case "get", "post", "put", "delete":
		if idx := strings.IndexByte(seq.Model, '.'); idx > 0 {
			model := seq.Model[:idx]
			modelRefs[model] = true
			// DDL table name = lowercase plural of model name
			modelRefs[strings.ToLower(inflection.Plural(model))] = true
		}
	case "response":
		populateResponseFields(g, funcName, seq)
	}
}
