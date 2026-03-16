//ff:func feature=gen-gogin type=generator control=sequence topic=interface-derive
//ff:what creates model/include_helpers.go with shared utility functions

package gogin

import (
	"os"
	"path/filepath"
	"strings"
)

// generateIncludeHelpersFile creates model/include_helpers.go with shared utility functions.
func generateIncludeHelpersFile(modelDir string) error {
	var b strings.Builder
	b.WriteString("package model\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString(")\n\n")

	b.WriteString("func collectInt64s(ids map[int64]bool) []int64 {\n")
	b.WriteString("\tkeys := make([]int64, 0, len(ids))\n")
	b.WriteString("\tfor k := range ids {\n")
	b.WriteString("\t\tkeys = append(keys, k)\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn keys\n")
	b.WriteString("}\n\n")

	b.WriteString("func buildPlaceholders(n int) string {\n")
	b.WriteString("\tps := make([]string, n)\n")
	b.WriteString("\tfor i := range ps {\n")
	b.WriteString("\t\tps[i] = fmt.Sprintf(\"$%d\", i+1)\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn strings.Join(ps, \", \")\n")
	b.WriteString("}\n\n")

	b.WriteString("func int64sToArgs(keys []int64) []interface{} {\n")
	b.WriteString("\targs := make([]interface{}, len(keys))\n")
	b.WriteString("\tfor i, k := range keys {\n")
	b.WriteString("\t\targs[i] = k\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn args\n")
	b.WriteString("}\n")

	return os.WriteFile(filepath.Join(modelDir, "include_helpers.go"), []byte(b.String()), 0644)
}
