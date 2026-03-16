//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=interface-derive
//ff:what generates a scan helper function for a model

package gogin

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
)

// generateScanFunc generates a scan helper function for a model.
func generateScanFunc(b *strings.Builder, modelName string, table *ddlTable) {
	lowerName := strcase.ToGoCamel(modelName)
	varName := string(lowerName[0])

	b.WriteString(fmt.Sprintf("func scan%s(row interface{ Scan(...interface{}) error }) (*%s, error) {\n", modelName, modelName))
	b.WriteString(fmt.Sprintf("\tvar %s %s\n", varName, modelName))

	scanFields := make([]string, len(table.Columns))
	for i, col := range table.Columns {
		scanFields[i] = fmt.Sprintf("&%s.%s", varName, col.GoName)
	}

	b.WriteString(fmt.Sprintf("\terr := row.Scan(%s)\n", strings.Join(scanFields, ", ")))
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\treturn nil, err\n")
	b.WriteString("\t}\n")
	b.WriteString(fmt.Sprintf("\treturn &%s, nil\n", varName))
	b.WriteString("}\n")
}
