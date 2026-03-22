//ff:func feature=gen-gogin type=generator control=selection topic=interface-derive
//ff:what writes a single method implementation based on the interface signature

package gogin

import (
	"fmt"
	"strings"
)

// generateMethodFromIface writes a single method implementation based on the interface signature.
func generateMethodFromIface(b *strings.Builder, implName, modelName string, m ifaceMethod, query *sqlcQuery, seqType string, table *ddlTable, includes []includeMapping, cursorSpecs map[string]string) {
	// WithTx special case: return new impl with tx set.
	if m.Name == "WithTx" {
		b.WriteString(fmt.Sprintf("func (m *%s) WithTx(tx *sql.Tx) %sModel {\n", implName, modelName))
		b.WriteString(fmt.Sprintf("\treturn &%s{db: m.db, tx: tx}\n", implName))
		b.WriteString("}\n")
		return
	}

	sqlStr := "-- TODO: " + m.Name
	if query != nil && query.SQL != "" {
		sqlStr = query.SQL
	}

	callArgNames := reorderCallArgs(m, query)
	callArgs := formatCallArgs(callArgNames)

	hasQueryOpts := hasQueryOptsParam(m)
	isList := isListMethod(m.Name) && hasQueryOpts
	isPageReturn := strings.Contains(m.ReturnSig, "pagination.Page[")
	isCursorReturn := strings.Contains(m.ReturnSig, "pagination.Cursor[")
	isFind := strings.HasPrefix(m.Name, "Find")
	isSliceReturn := strings.Contains(m.ReturnSig, "[]")

	switch {
	case isList && isCursorReturn:
		writeCursorPaginationMethod(b, implName, modelName, m, query, table, includes, callArgNames, callArgs, cursorSpecs)

	case isList:
		writeOffsetPaginationMethod(b, implName, modelName, m, query, table, includes, callArgNames, callArgs, isPageReturn)

	case isSliceReturn:
		writeSliceReturnMethod(b, implName, modelName, m, sqlStr, callArgs)

	case isFind || seqType == "get":
		writeFindMethod(b, implName, modelName, m, sqlStr, callArgs)

	case seqType == "post":
		b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
		b.WriteString(fmt.Sprintf("\trow := m.conn().QueryRowContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
		b.WriteString(fmt.Sprintf("\treturn scan%s(row)\n", modelName))
		b.WriteString("}\n")

	case seqType == "put" || seqType == "delete":
		b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
		b.WriteString(fmt.Sprintf("\t_, err := m.conn().ExecContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
		b.WriteString("\treturn err\n")
		b.WriteString("}\n")

	default:
		if query != nil && query.Cardinality == "one" {
			b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
			b.WriteString(fmt.Sprintf("\trow := m.conn().QueryRowContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
			b.WriteString(fmt.Sprintf("\treturn scan%s(row)\n", modelName))
			b.WriteString("}\n")
		} else {
			b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
			b.WriteString(fmt.Sprintf("\t_, err := m.conn().ExecContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
			b.WriteString("\treturn err\n")
			b.WriteString("}\n")
		}
	}
}
