//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=interface-derive
//ff:what offset-based pagination 메서드 구현을 생성한다

package gogin

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
)

// writeOffsetPaginationMethod writes an offset-based pagination method implementation.
func writeOffsetPaginationMethod(b *strings.Builder, implName, modelName string, m ifaceMethod, query *sqlcQuery, table *ddlTable, includes []includeMapping, callArgNames []string, callArgs string, isPageReturn bool) {
	baseWhere := ""
	baseArgCount := 0
	if query != nil && query.SQL != "" {
		baseWhere, baseArgCount = extractBaseWhere(query.SQL)
	}

	tableName := ""
	if table != nil {
		tableName = table.TableName
	}

	b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))

	if len(callArgNames) > 0 {
		b.WriteString(fmt.Sprintf("\tbaseArgs := []interface{}{%s}\n", strings.Join(callArgNames, ", ")))
	}

	b.WriteString(fmt.Sprintf("\tcountSQL, countArgs := BuildCountQuery(%q, %q, %d, opts)\n", tableName, baseWhere, baseArgCount))
	if len(callArgNames) > 0 {
		b.WriteString("\tcountArgs = append(baseArgs, countArgs...)\n")
	}
	b.WriteString("\tvar total int64\n")
	b.WriteString("\tif err := m.conn().QueryRowContext(context.Background(), countSQL, countArgs...).Scan(&total); err != nil {\n")
	b.WriteString("\t\treturn nil, err\n")
	b.WriteString("\t}\n\n")

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
	for _, inc := range includes {
		helperName := "include" + strcase.ToGoPascal(inc.IncludeName)
		b.WriteString(fmt.Sprintf("\tif err := m.%s(items); err != nil {\n", helperName))
		b.WriteString("\t\treturn nil, err\n")
		b.WriteString("\t}\n")
	}
	if isPageReturn {
		b.WriteString(fmt.Sprintf("\treturn &pagination.Page[%s]{Items: items, Total: total}, nil\n", modelName))
	} else {
		b.WriteString("\treturn items, total, nil\n")
	}
	b.WriteString("}\n")
}
