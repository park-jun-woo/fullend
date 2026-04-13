//ff:func feature=gen-gogin type=generator control=sequence topic=interface-derive
//ff:what generates a forward FK include helper method for a model

package gogin

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
)

// generateIncludeHelper generates a forward FK include helper method for a model.
func generateIncludeHelper(b *strings.Builder, implName, modelName string, inc includeMapping) {
	helperName := "include" + strcase.ToGoPascal(inc.IncludeName)
	fkGoName := snakeToGo(inc.FKColumn)

	b.WriteString(fmt.Sprintf("func (m *%s) %s(items []%s) error {\n", implName, helperName, modelName))
	b.WriteString("\tids := make(map[int64]bool)\n")
	b.WriteString("\tfor _, item := range items {\n")
	b.WriteString(fmt.Sprintf("\t\tids[item.%s] = true\n", fkGoName))
	b.WriteString("\t}\n")
	b.WriteString("\tif len(ids) == 0 {\n")
	b.WriteString("\t\treturn nil\n")
	b.WriteString("\t}\n")
	b.WriteString("\tkeys := collectInt64s(ids)\n")
	b.WriteString("\tplaceholders := buildPlaceholders(len(keys))\n")
	b.WriteString("\targs := int64sToArgs(keys)\n")
	b.WriteString(fmt.Sprintf("\trows, err := m.conn().QueryContext(context.Background(),\n\t\t\"SELECT * FROM %s WHERE id IN (\"+placeholders+\")\", args...)\n", inc.TargetTable))
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\treturn err\n")
	b.WriteString("\t}\n")
	b.WriteString("\tdefer rows.Close()\n")
	b.WriteString(fmt.Sprintf("\tlookup := make(map[int64]*%s)\n", inc.TargetModel))
	b.WriteString("\tfor rows.Next() {\n")
	b.WriteString(fmt.Sprintf("\t\tv, err := scan%s(rows)\n", inc.TargetModel))
	b.WriteString("\t\tif err != nil {\n")
	b.WriteString("\t\t\treturn err\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t\tlookup[v.ID] = v\n")
	b.WriteString("\t}\n")
	b.WriteString("\tfor i := range items {\n")
	b.WriteString(fmt.Sprintf("\t\titems[i].%s = lookup[items[i].%s]\n", inc.FieldName, fkGoName))
	b.WriteString("\t}\n")
	b.WriteString("\treturn nil\n")
	b.WriteString("}\n")
}
