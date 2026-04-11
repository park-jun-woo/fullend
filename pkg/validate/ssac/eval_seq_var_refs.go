//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what evalSeqVarRefs — 시퀀스의 Args/Inputs에서 변수 참조를 추출하여 VarDeclared 평가
package ssac

import (
	"strings"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func evalSeqVarRefs(graph *toulmin.Graph, g *rule.Ground, file, funcName string, seqIdx int, seq parsessac.Sequence) []validate.ValidationError {
	var refs []string
	for _, arg := range seq.Args {
		if arg.Source != "" && arg.Source != "request" && arg.Source != "currentUser" && arg.Source != "query" && arg.Source != "message" {
			refs = append(refs, arg.Source)
		}
	}
	for _, val := range seq.Inputs {
		if strings.HasPrefix(val, `"`) || parsessac.IsLiteral(val) {
			continue
		}
		ref := strings.SplitN(val, ".", 2)[0]
		if ref != "" && ref != "request" && ref != "currentUser" && ref != "query" && ref != "message" {
			refs = append(refs, ref)
		}
	}
	// S-29: Target variable ref
	if seq.Target != "" {
		ref := strings.SplitN(seq.Target, ".", 2)[0]
		if ref != "" && ref != "request" && ref != "currentUser" && ref != "query" && ref != "message" {
			refs = append(refs, ref)
		}
	}
	// S-30: response Fields variable refs
	for _, val := range seq.Fields {
		ref := strings.SplitN(val, ".", 2)[0]
		if ref != "" && ref != "request" && ref != "currentUser" && ref != "query" && ref != "message" {
			refs = append(refs, ref)
		}
	}
	var errs []validate.ValidationError
	for _, ref := range refs {
		ctx := toulmin.NewContext()
		ctx.Set("ground", g)
		ctx.Set("claim", ref)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toValidationErrors(results, file, funcName, seqIdx)...)
	}
	return errs
}
