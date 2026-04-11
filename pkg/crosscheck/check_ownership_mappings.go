//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkOwnershipMappings — 개별 ownership 매핑의 table/column 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/rego"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkOwnershipMappings(tableGraph *toulmin.Graph, g *rule.Ground, oms []rego.OwnershipMapping) []CrossError {
	var errs []CrossError
	for _, om := range oms {
		errs = append(errs, evalRef(tableGraph, g, om.Table, om.Resource)...)
		errs = append(errs, evalColumnRef(g, om.Table, om.Column, "X-32", om.Resource)...)
		if om.JoinTable != "" {
			errs = append(errs, evalRef(tableGraph, g, om.JoinTable, om.Resource+" via")...)
			errs = append(errs, evalColumnRef(g, om.JoinTable, om.JoinFK, "X-34", om.Resource+" via")...)
		}
	}
	return errs
}
