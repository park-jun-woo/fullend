//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=interface-derive
//ff:what writes method implementations and include helpers for a model file

package gogin

import "strings"

// writeModelMethods writes method implementations and include helpers.
func writeModelMethods(b *strings.Builder, modelName string, methods []ifaceMethod, table *ddlTable, queries map[string]sqlcQuery, seqTypes map[string]string, includes []includeMapping, cursorSpecs map[string]string, implName string) {
	for _, method := range methods {
		b.WriteString("\n")
		if table != nil {
			d := buildMethodDirective(table, method)
			b.WriteString(d.String() + "\n")
		}
		query := queries[method.Name]
		seqType := seqTypes[method.Name]
		generateMethodFromIface(b, implName, modelName, method, &query, seqType, table, includes, cursorSpecs)
	}

	for _, inc := range includes {
		b.WriteString("\n")
		generateIncludeHelper(b, implName, modelName, inc)
	}
}
