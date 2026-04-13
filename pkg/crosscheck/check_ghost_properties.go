//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkGhostProperties — OpenAPI property가 DDL에 없는 유령 property 감지 (X-9, X-10)
package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkGhostProperties(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil || len(g.Tables) == 0 {
		return nil
	}
	var errs []CrossError
	for opID, fields := range g.Schemas {
		if !strings.HasPrefix(opID, "OpenAPI.response.resolved.") {
			continue
		}
		op := opID[len("OpenAPI.response.resolved."):]
		table := guessTableFromOp(op)
		t, ok := g.Tables[table]
		if !ok || len(t.Columns) == 0 {
			continue
		}
		errs = append(errs, checkGhostFields(op, table, fields, t.Columns)...)
	}
	return errs
}
