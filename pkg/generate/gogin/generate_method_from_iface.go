//ff:func feature=gen-gogin type=generator control=selection topic=interface-derive
//ff:what dispatches a method implementation based on DecideMethodPattern

package gogin

import (
	"fmt"
	"strings"
)

// generateMethodFromIface dispatches to the matching Pattern handler.
// Decision logic lives in DecideMethodPattern; this function is a thin dispatcher.
func generateMethodFromIface(b *strings.Builder, implName, modelName string, m ifaceMethod, query *sqlcQuery, seqType string, table *ddlTable, includes []includeMapping, cursorSpecs map[string]string) {
	facts := NewMethodFacts(m, query, seqType)
	pattern := DecideMethodPattern(facts)

	if pattern == PatternSkip {
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
	isPageReturn := strings.Contains(m.ReturnSig, "pagination.Page[")

	switch pattern {
	case PatternCursorPagination:
		writeCursorPaginationMethod(b, implName, modelName, m, query, table, includes, callArgNames, callArgs, cursorSpecs)

	case PatternOffsetPagination:
		writeOffsetPaginationMethod(b, implName, modelName, m, query, table, includes, callArgNames, callArgs, isPageReturn)

	case PatternSliceReturn:
		writeSliceReturnMethod(b, implName, modelName, m, sqlStr, callArgs)

	case PatternFind:
		writeFindMethod(b, implName, modelName, m, sqlStr, callArgs)

	case PatternCreate:
		b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
		b.WriteString(fmt.Sprintf("\trow := m.conn().QueryRowContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
		b.WriteString(fmt.Sprintf("\treturn scan%s(row)\n", modelName))
		b.WriteString("}\n")

	case PatternUpdateDelete:
		b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
		b.WriteString(fmt.Sprintf("\t_, err := m.conn().ExecContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
		b.WriteString("\treturn err\n")
		b.WriteString("}\n")

	case PatternFallbackOne:
		b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
		b.WriteString(fmt.Sprintf("\trow := m.conn().QueryRowContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
		b.WriteString(fmt.Sprintf("\treturn scan%s(row)\n", modelName))
		b.WriteString("}\n")

	case PatternFallbackExec:
		b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
		b.WriteString(fmt.Sprintf("\t_, err := m.conn().ExecContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
		b.WriteString("\treturn err\n")
		b.WriteString("}\n")
	}
}
