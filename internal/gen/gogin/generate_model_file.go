//ff:func feature=gen-gogin type=generator control=iteration dimension=2
//ff:what creates model/{model}.go with the implementation struct using *sql.DB

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

)

// generateModelFile creates model/{model}.go with the implementation struct using *sql.DB.
func generateModelFile(modelDir string, modelName string, methods []ifaceMethod, table *ddlTable, queries map[string]sqlcQuery, seqTypes map[string]string, includes []includeMapping, cursorSpecs map[string]string) error {
	var b strings.Builder
	lowerName := strings.ToLower(modelName)
	implName := lowerName + "ModelImpl"

	b.WriteString("package model\n\n")

	// Check if any method returns pagination.Page or pagination.Cursor.
	needsPagination := false
	needsCursor := false
	for _, method := range methods {
		if strings.Contains(method.ReturnSig, "pagination.Page[") || strings.Contains(method.ReturnSig, "pagination.Cursor[") {
			needsPagination = true
		}
		if strings.Contains(method.ReturnSig, "pagination.Cursor[") {
			needsCursor = true
		}
	}

	needsJSON := false
	for _, method := range methods {
		for _, p := range method.Params {
			if p.Type == "json.RawMessage" {
				needsJSON = true
				break
			}
		}
		if needsJSON {
			break
		}
	}

	b.WriteString("import (\n")
	b.WriteString("\t\"context\"\n")
	b.WriteString("\t\"database/sql\"\n")
	if needsJSON {
		b.WriteString("\t\"encoding/json\"\n")
	}
	if needsCursor {
		b.WriteString("\t\"fmt\"\n")
	}
	if needsPagination {
		b.WriteString("\n\t\"github.com/geul-org/fullend/pkg/pagination\"\n")
	}
	b.WriteString(")\n\n")

	// Struct definition.
	b.WriteString(fmt.Sprintf("type %s struct {\n", implName))
	b.WriteString("\tdb *sql.DB\n")
	b.WriteString("\ttx *sql.Tx\n")
	b.WriteString("}\n\n")

	// conn() helper — returns tx if set, otherwise db.
	b.WriteString(fmt.Sprintf("func (m *%s) conn() interface {\n", implName))
	b.WriteString("\tExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)\n")
	b.WriteString("\tQueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)\n")
	b.WriteString("\tQueryRowContext(ctx context.Context, query string, args ...any) *sql.Row\n")
	b.WriteString("} {\n")
	b.WriteString("\tif m.tx != nil {\n")
	b.WriteString("\t\treturn m.tx\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn m.db\n")
	b.WriteString("}\n\n")

	// Constructor.
	b.WriteString(fmt.Sprintf("func New%sModel(db *sql.DB) %sModel {\n", modelName, modelName))
	b.WriteString(fmt.Sprintf("\treturn &%s{db: db}\n", implName))
	b.WriteString("}\n")

	// Scan helper.
	if table != nil {
		b.WriteString("\n")
		generateScanFunc(&b, modelName, table)
	}

	writeModelMethods(&b, modelName, methods, table, queries, seqTypes, includes, cursorSpecs, implName)

	fileName := lowerName + ".go"
	return os.WriteFile(filepath.Join(modelDir, fileName), []byte(b.String()), 0644)
}
