//ff:func feature=gen-gogin type=generator control=sequence topic=interface-derive
//ff:what pagination 없는 다건 조회 메서드 구현을 생성한다

package gogin

import (
	"fmt"
	"strings"
)

// writeSliceReturnMethod writes a multi-row query method without pagination.
func writeSliceReturnMethod(b *strings.Builder, implName, modelName string, m ifaceMethod, sqlStr, callArgs string) {
	b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
	b.WriteString(fmt.Sprintf("\trows, err := m.conn().QueryContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
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
	b.WriteString("\treturn items, nil\n")
	b.WriteString("}\n")
}
