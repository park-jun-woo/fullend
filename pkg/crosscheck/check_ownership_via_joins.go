//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkOwnershipViaJoins — 단일 정책의 @ownership join table DDL 존재 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/rego"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkOwnershipViaJoins(tableGraph *toulmin.Graph, g *rule.Ground, ownerships []rego.OwnershipMapping) []CrossError {
	var errs []CrossError
	for _, om := range ownerships {
		if om.JoinTable != "" {
			errs = append(errs, evalRef(tableGraph, g, om.JoinTable, om.Resource+" via")...)
		}
	}
	return errs
}
