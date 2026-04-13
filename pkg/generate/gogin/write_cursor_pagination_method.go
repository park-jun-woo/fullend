//ff:func feature=gen-gogin type=generator control=sequence topic=interface-derive
//ff:what cursor-based pagination 메서드 구현을 생성한다

package gogin

import (
	"fmt"
	"strings"
)

// writeCursorPaginationMethod writes a cursor-based pagination method implementation.
func writeCursorPaginationMethod(b *strings.Builder, implName, modelName string, m ifaceMethod, query *sqlcQuery, table *ddlTable, includes []includeMapping, callArgNames []string, callArgs string, cursorSpecs map[string]string) {
	baseWhere := ""
	baseArgCount := 0
	if query != nil && query.SQL != "" {
		baseWhere, baseArgCount = extractBaseWhere(query.SQL)
	}

	tableName := ""
	if table != nil {
		tableName = table.TableName
	}

	cursorField := lookupCursorField(cursorSpecs, m.Name)

	b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))

	if len(callArgNames) > 0 {
		b.WriteString(fmt.Sprintf("\tbaseArgs := []interface{}{%s}\n", strings.Join(callArgNames, ", ")))
	}

	b.WriteString("\trequestedLimit := opts.Limit\n")
	b.WriteString("\topts.Limit = requestedLimit + 1\n\n")

	b.WriteString(fmt.Sprintf("\tselectSQL, selectArgs := BuildSelectQuery(%q, %q, %d, opts)\n", tableName, baseWhere, baseArgCount))
	if len(callArgNames) > 0 {
		b.WriteString("\tselectArgs = append(baseArgs, selectArgs...)\n")
	}
	b.WriteString("\trows, err := m.conn().QueryContext(context.Background(), selectSQL, selectArgs...)\n")
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\treturn nil, err\n")
	b.WriteString("\t}\n")
	b.WriteString("\tdefer rows.Close()\n")
	b.WriteString(fmt.Sprintf("\titems := make([]%s, 0)\n", modelName))
	b.WriteString("\tfor rows.Next() {\n")
	b.WriteString(fmt.Sprintf("\t\tv, err := scan%s(rows)\n", modelName))
	b.WriteString("\t\tif err != nil {\n")
	b.WriteString("\t\t\treturn nil, err\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t\titems = append(items, *v)\n")
	b.WriteString("\t}\n")
	b.WriteString("\tif err := rows.Err(); err != nil {\n")
	b.WriteString("\t\treturn nil, err\n")
	b.WriteString("\t}\n")

	writeIncludeLoads(b, includes)

	b.WriteString("\thasNext := len(items) > requestedLimit\n")
	b.WriteString("\tvar nextCursor string\n")
	b.WriteString("\tif hasNext {\n")
	b.WriteString("\t\titems = items[:requestedLimit]\n")
	b.WriteString("\t}\n")
	b.WriteString("\tif len(items) > 0 {\n")
	b.WriteString(fmt.Sprintf("\t\tnextCursor = fmt.Sprintf(\"%%v\", items[len(items)-1].%s)\n", cursorField))
	b.WriteString("\t}\n")
	b.WriteString(fmt.Sprintf("\treturn &pagination.Cursor[%s]{Items: items, NextCursor: nextCursor, HasNext: hasNext}, nil\n", modelName))
	b.WriteString("}\n")
}
