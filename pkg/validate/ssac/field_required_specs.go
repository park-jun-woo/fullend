//ff:func feature=rule type=util control=selection
//ff:what fieldRequiredSpecs — 시퀀스 타입별 FieldRequired Spec 목록 반환
package ssac

import (
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func fieldRequiredSpecs(seqType string) []toulmin.Spec {
	switch seqType {
	case "get":
		return []toulmin.Spec{
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-1", Level: "ERROR", Message: "@get requires Model"}, SeqType: "get", Field: "Model", Required: true},
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-2", Level: "ERROR", Message: "@get requires Result"}, SeqType: "get", Field: "Result", Required: true},
		}
	case "post":
		return []toulmin.Spec{
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-3", Level: "ERROR", Message: "@post requires Model"}, SeqType: "post", Field: "Model", Required: true},
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-4", Level: "ERROR", Message: "@post requires Result"}, SeqType: "post", Field: "Result", Required: true},
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-5", Level: "ERROR", Message: "@post requires Inputs"}, SeqType: "post", Field: "Args", Required: true},
		}
	case "put":
		return []toulmin.Spec{
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-6", Level: "ERROR", Message: "@put requires Model"}, SeqType: "put", Field: "Model", Required: true},
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-7", Level: "ERROR", Message: "@put must not have Result"}, SeqType: "put", Field: "Result", Required: false},
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-8", Level: "ERROR", Message: "@put requires Inputs"}, SeqType: "put", Field: "Args", Required: true},
		}
	case "delete":
		return []toulmin.Spec{
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-9", Level: "ERROR", Message: "@delete requires Model"}, SeqType: "delete", Field: "Model", Required: true},
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-10", Level: "ERROR", Message: "@delete must not have Result"}, SeqType: "delete", Field: "Result", Required: false},
		}
	case "empty":
		return []toulmin.Spec{
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-12", Level: "ERROR", Message: "@empty requires Target"}, SeqType: "empty", Field: "Target", Required: true},
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-13", Level: "ERROR", Message: "@empty requires Message"}, SeqType: "empty", Field: "Message", Required: true},
		}
	case "state":
		return []toulmin.Spec{
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-14", Level: "ERROR", Message: "@state requires DiagramID"}, SeqType: "state", Field: "DiagramID", Required: true},
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-15", Level: "ERROR", Message: "@state requires Inputs"}, SeqType: "state", Field: "Inputs", Required: true},
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-16", Level: "ERROR", Message: "@state requires Transition"}, SeqType: "state", Field: "Transition", Required: true},
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-17", Level: "ERROR", Message: "@state requires Message"}, SeqType: "state", Field: "Message", Required: true},
		}
	case "auth":
		return []toulmin.Spec{
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-18", Level: "ERROR", Message: "@auth requires Action"}, SeqType: "auth", Field: "Action", Required: true},
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-19", Level: "ERROR", Message: "@auth requires Resource"}, SeqType: "auth", Field: "Resource", Required: true},
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-20", Level: "ERROR", Message: "@auth requires Message"}, SeqType: "auth", Field: "Message", Required: true},
		}
	case "call":
		return []toulmin.Spec{
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-21", Level: "ERROR", Message: "@call requires Model (package.Func)"}, SeqType: "call", Field: "Model", Required: true},
		}
	case "publish":
		return []toulmin.Spec{
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-23", Level: "ERROR", Message: "@publish requires Topic"}, SeqType: "publish", Field: "Topic", Required: true},
			&rule.FieldRequiredSpec{BaseSpec: rule.BaseSpec{Rule: "S-24", Level: "ERROR", Message: "@publish requires Payload"}, SeqType: "publish", Field: "Payload", Required: true},
		}
	default:
		return nil
	}
}
