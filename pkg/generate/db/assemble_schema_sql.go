//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=ddl
//ff:what assembleSchemaSQL — 위상순서대로 CREATE TABLE/INDEX + INSERT seed 병합

package db

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/ddl"
)

// assembleSchemaSQL concatenates per-table DDL content (file-level) in order,
// followed by auto-generated seeds. Returns (sql, totalSeeds).
func assembleSchemaSQL(order []string, tableByName map[string]*ddl.Table, specsDir, autoSeeds string) (string, int) {
	var sb strings.Builder
	totalSeeds := 0
	for _, name := range order {
		t := tableByName[name]
		if t == nil {
			continue
		}
		fileContent := readDDLFileForTable(specsDir, name)
		if fileContent == "" {
			continue
		}
		sb.WriteString(fmt.Sprintf("-- ---- %s ----\n", name))
		sb.WriteString(strings.TrimSpace(fileContent))
		sb.WriteString("\n\n")
		totalSeeds += len(t.Seeds)
	}
	if autoSeeds != "" {
		sb.WriteString("-- ---- auto-generated nobody seeds ----\n")
		sb.WriteString(autoSeeds)
		sb.WriteString("\n")
	}
	return sb.String(), totalSeeds
}
