//ff:func feature=gen-gogin type=generator control=sequence topic=interface-derive
//ff:what 단건 조회(Find/Get) 메서드 구현을 생성한다

package gogin

import (
	"fmt"
	"strings"
)

// writeFindMethod writes a single-row query method (Find/Get).
func writeFindMethod(b *strings.Builder, implName, modelName string, m ifaceMethod, sqlStr, callArgs string) {
	b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
	b.WriteString(fmt.Sprintf("\trow := m.conn().QueryRowContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
	b.WriteString(fmt.Sprintf("\tv, err := scan%s(row)\n", modelName))
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\tif err == sql.ErrNoRows {\n")
	b.WriteString("\t\t\treturn nil, nil\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t\treturn nil, err\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn v, nil\n")
	b.WriteString("}\n")
}
