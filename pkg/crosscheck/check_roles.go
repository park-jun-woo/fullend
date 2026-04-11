//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkRoles — Rego roles ↔ Config roles, DDL CHECK 교차 검증 (X-63~X-65)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkRoles(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.ParsedPolicies) == 0 || fs.Manifest == nil {
		return nil
	}
	var errs []CrossError

	graph63 := toulmin.NewGraph("rego-config-roles")
	graph63.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-63", Level: "ERROR", Message: "Rego role not in fullend.yaml auth.roles"},
		LookupKey: "Config.roles",
	})
	for rv := range g.Lookup["Rego.roles"] {
		errs = append(errs, evalRef(graph63, g, rv, rv)...)
	}

	graph64 := toulmin.NewGraph("config-rego-roles")
	graph64.Rule(rule.CoverageCheck).With(&rule.CoverageCheckSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-64", Level: "WARNING", Message: "Config role not referenced in Rego"},
		LookupKey: "Rego.roles",
	})
	for rv := range g.Lookup["Config.roles"] {
		errs = append(errs, evalRef(graph64, g, rv, rv)...)
	}

	return errs
}
