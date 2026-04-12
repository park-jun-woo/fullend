//ff:func feature=rule type=generator control=sequence
//ff:what evaluateSeq — 단일 시퀀스를 Graph로 평가하여 Trace Pattern 반환
package backend

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/generate/trace"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func evaluateSeq(graph *toulmin.Graph, seq parsessac.Sequence, funcName string, fkRef bool) trace.Pattern {
	claim := SeqClaim{Type: seq.Type, Seq: seq, FuncName: funcName, FKRef: fkRef}
	ctx := toulmin.NewContext()
	ctx.Set("claim", claim)
	results, _ := graph.Evaluate(ctx, toulmin.EvalOption{Trace: true})
	return trace.BuildPattern(results)
}
