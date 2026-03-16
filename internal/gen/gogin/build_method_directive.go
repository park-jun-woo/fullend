//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=output
//ff:what builds a contract directive for a model method from DDL table info

package gogin

import "github.com/geul-org/fullend/internal/contract"

// buildMethodDirective builds a contract directive for a model method.
func buildMethodDirective(table *ddlTable, method ifaceMethod) *contract.Directive {
	ssotPath := "db/" + table.TableName + ".sql"
	params := make([]string, len(method.Params))
	for i, p := range method.Params {
		params[i] = p.Type
	}
	returns := parseReturnTypes(method.ReturnSig)
	hash := contract.HashModelMethod(method.Name, params, returns)
	return &contract.Directive{Ownership: "gen", SSOT: ssotPath, Contract: hash}
}
