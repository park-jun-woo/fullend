//ff:func feature=gen-gogin type=generator control=iteration
//ff:what creates model/types.go with struct definitions from DDL columns

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// generateTypesFile creates model/types.go with struct definitions from DDL columns.
func generateTypesFile(modelDir string, models []string, tables map[string]*ddlTable, includesByModel map[string][]includeMapping) error {
	var b strings.Builder
	b.WriteString("package model\n\n")

	// Determine if we need time/json imports.
	needsTime := false
	needsJSON := false
	for _, m := range models {
		t := tables[m]
		if t == nil {
			continue
		}
		for _, col := range t.Columns {
			if col.GoType == "time.Time" {
				needsTime = true
			}
			if col.GoType == "json.RawMessage" {
				needsJSON = true
			}
		}
		if needsTime && needsJSON {
			break
		}
	}

	if needsJSON || needsTime {
		var imports []string
		if needsJSON {
			imports = append(imports, "\"encoding/json\"")
		}
		if needsTime {
			imports = append(imports, "\"time\"")
		}
		b.WriteString("import (\n")
		for _, imp := range imports {
			b.WriteString("\t" + imp + "\n")
		}
		b.WriteString(")\n\n")
	}

	for i, m := range models {
		t := tables[m]
		if t == nil {
			continue
		}
		b.WriteString(fmt.Sprintf("type %s struct {\n", m))
		for _, col := range t.Columns {
			jsonTag := col.Name
			if col.Sensitive {
				jsonTag = "-"
			}
			b.WriteString(fmt.Sprintf("\t%-12s %s `json:\"%s\"`\n", col.GoName, col.GoType, jsonTag))
		}
		if includes, ok := includesByModel[m]; ok && len(includes) > 0 {
			b.WriteString("\n\t// Include fields\n")
			for _, inc := range includes {
				jsonTag := lcFirst(inc.FieldName)
				b.WriteString(fmt.Sprintf("\t%-12s %s `json:\"%s,omitempty\"`\n", inc.FieldName, inc.FieldType, jsonTag))
			}
		}
		b.WriteString("}\n")
		if i < len(models)-1 {
			b.WriteString("\n")
		}
	}

	return os.WriteFile(filepath.Join(modelDir, "types.go"), []byte(b.String()), 0644)
}
