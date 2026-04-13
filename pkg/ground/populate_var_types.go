//ff:func feature=rule type=loader control=iteration dimension=1
//ff:what populateVarTypes — SSaC 시퀀스 result에서 변수→타입 매핑 구축
package ground

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateVarTypes(g *rule.Ground, fs *fullend.Fullstack) {
	for _, fn := range fs.ServiceFuncs {
		populateVarTypesSeqs(g, fn.Name, fn.Sequences)
	}
}
