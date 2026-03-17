//ff:func feature=crosscheck type=rule control=sequence topic=policy-check
//ff:what 정책과 SSaC/DDL/States 간 교차 검증 실행
package crosscheck

import (
	"github.com/park-jun-woo/fullend/internal/policy"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
	"github.com/park-jun-woo/fullend/internal/statemachine"
)

// CheckPolicy validates policy against SSaC, DDL, and States.
func CheckPolicy(policies []*policy.Policy, funcs []ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable, diagrams []*statemachine.StateDiagram) []CrossError {
	var errs []CrossError

	allPairs, ownerResources, allOwnerships := mergePolicies(policies)
	ssacPairs := buildSSaCAuthPairs(funcs)

	errs = append(errs, checkSSaCPairsCoverage(ssacPairs, allPairs)...)
	errs = append(errs, checkRegoPairsCoverage(allPairs, ssacPairs)...)
	errs = append(errs, checkOwnershipAnnotations(ownerResources, allOwnerships)...)

	if st != nil {
		errs = append(errs, checkOwnershipDDL(allOwnerships, st)...)
	}

	return errs
}
