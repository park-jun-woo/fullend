//ff:func feature=orchestrator type=util control=iteration
//ff:what traceDDL finds DDL tables referenced by SSaC sequences.

package orchestrator

import (
	"strings"

	"github.com/jinzhu/inflection"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

func traceDDL(sf *ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable, specsDir string) []ChainLink {
	tables := map[string]bool{}
	for _, seq := range sf.Sequences {
		if seq.Model == "" || seq.Type == "call" || seq.Type == "response" {
			continue
		}
		parts := strings.SplitN(seq.Model, ".", 2)
		if len(parts) < 2 {
			continue
		}
		tableName := inflection.Plural(toSnakeCase(parts[0]))
		if _, ok := st.DDLTables[tableName]; ok {
			tables[tableName] = true
		}
	}

	var links []ChainLink
	sortedTables := sortedStringKeys(tables)
	for _, table := range sortedTables {
		// Find the DDL file.
		relPath, line := findDDLTable(table, specsDir)
		links = append(links, ChainLink{
			Kind:    "DDL",
			File:    relPath,
			Line:    line,
			Summary: "CREATE TABLE " + table,
		})
	}
	return links
}
