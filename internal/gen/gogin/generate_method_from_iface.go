//ff:func feature=gen-gogin type=generator
//ff:what writes a single method implementation based on the interface signature

package gogin

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
)

// generateMethodFromIface writes a single method implementation based on the interface signature.
func generateMethodFromIface(b *strings.Builder, implName, modelName string, m ifaceMethod, query *sqlcQuery, seqType string, table *ddlTable, includes []includeMapping, cursorSpecs map[string]string) {
	// WithTx special case: return new impl with tx set.
	if m.Name == "WithTx" {
		b.WriteString(fmt.Sprintf("func (m *%s) WithTx(tx *sql.Tx) %sModel {\n", implName, modelName))
		b.WriteString(fmt.Sprintf("\treturn &%s{db: m.db, tx: tx}\n", implName))
		b.WriteString("}\n")
		return
	}

	sqlStr := "-- TODO: " + m.Name
	if query != nil && query.SQL != "" {
		sqlStr = query.SQL
	}

	// Build call args from interface params (excluding QueryOpts params).
	// For INSERT/UPDATE, reorder args to match SQL column order from sqlcQuery.Columns.
	var callArgNames []string
	if query != nil && len(query.Columns) > 0 && len(m.Params) > 0 {
		// Build param lookup: goName (lowercase first) → param name.
		paramByCol := make(map[string]string) // sql_column → param name
		for _, p := range m.Params {
			if p.Type == "QueryOpts" {
				continue
			}
			// Convert param name (camelCase) to snake_case for matching.
			snakeName := goToSnake(p.Name)
			paramByCol[snakeName] = p.Name
		}
		// Reorder: follow SQL column order, then append unmatched WHERE params.
		matched := make(map[string]bool)
		for _, col := range query.Columns {
			if paramName, ok := paramByCol[col]; ok {
				callArgNames = append(callArgNames, paramName)
				matched[paramName] = true
			}
		}
		// Append remaining params not matched by columns (e.g. WHERE id = $N).
		for _, p := range m.Params {
			if p.Type == "QueryOpts" {
				continue
			}
			if !matched[p.Name] {
				callArgNames = append(callArgNames, p.Name)
			}
		}
	} else {
		for _, p := range m.Params {
			if p.Type == "QueryOpts" {
				continue
			}
			callArgNames = append(callArgNames, p.Name)
		}
	}
	callArgs := ""
	if len(callArgNames) > 0 {
		callArgs = ",\n\t\t" + strings.Join(callArgNames, ", ")
	}

	// Determine method pattern from return signature and seq type.
	// A List method with QueryOpts uses dynamic SQL with pagination/sort/filter.
	// A List method without QueryOpts is a simple query returning a slice.
	hasQueryOpts := false
	for _, p := range m.Params {
		if p.Type == "QueryOpts" {
			hasQueryOpts = true
			break
		}
	}
	isList := isListMethod(m.Name) && hasQueryOpts
	isPageReturn := strings.Contains(m.ReturnSig, "pagination.Page[")
	isCursorReturn := strings.Contains(m.ReturnSig, "pagination.Cursor[")
	isFind := strings.HasPrefix(m.Name, "Find")

	// Check if return type is a slice (e.g. "[]Lesson" in "([]Lesson, error)").
	isSliceReturn := strings.Contains(m.ReturnSig, "[]")

	switch {
	case isList && isCursorReturn:
		// Cursor-based pagination: no COUNT, LIMIT+1 for hasNext detection.
		baseWhere := ""
		baseArgCount := 0
		if query != nil && query.SQL != "" {
			baseWhere, baseArgCount = extractBaseWhere(query.SQL)
		}

		tableName := ""
		if table != nil {
			tableName = table.TableName
		}

		// Determine cursor field name from cursorSpecs.
		cursorField := "ID"
		if cursorSpecs != nil {
			// Try to find operationId by matching method name to known cursor operations.
			for opID, field := range cursorSpecs {
				if opID == m.Name || strings.EqualFold(opID, m.Name) {
					cursorField = field
					break
				}
			}
		}

		b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))

		// Build base args from non-opts params.
		if len(callArgNames) > 0 {
			b.WriteString(fmt.Sprintf("\tbaseArgs := []interface{}{%s}\n", strings.Join(callArgNames, ", ")))
		}

		// Save requested limit, then bump to LIMIT+1 for hasNext detection.
		b.WriteString("\trequestedLimit := opts.Limit\n")
		b.WriteString("\topts.Limit = requestedLimit + 1\n\n")

		// Select query.
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

		// Include loading.
		for _, inc := range includes {
			helperName := "include" + strcase.ToGoPascal(inc.IncludeName)
			b.WriteString(fmt.Sprintf("\tif err := m.%s(items); err != nil {\n", helperName))
			b.WriteString("\t\treturn nil, err\n")
			b.WriteString("\t}\n")
		}

		// hasNext detection and cursor extraction.
		b.WriteString("\thasNext := len(items) > requestedLimit\n")
		b.WriteString("\tvar nextCursor string\n")
		b.WriteString("\tif hasNext {\n")
		b.WriteString("\t\titems = items[:requestedLimit]\n")
		b.WriteString("\t}\n")
		b.WriteString("\tif len(items) > 0 {\n")
		b.WriteString(fmt.Sprintf("\t\tnextCursor = fmt.Sprintf(\"%%v\", items[len(items)-1].%s)\n", cursorField))
		b.WriteString("\t}\n")
		b.WriteString(fmt.Sprintf("\treturn &pagination.Cursor[%s]{Items: items, NextCursor: nextCursor, HasNext: hasNext}, nil\n", modelName))
		b.WriteString("}\n")

	case isList:
		// Offset-based pagination: COUNT + SELECT returning *pagination.Page[T] or ([]T, int, error).
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

		// Build base args from non-opts params.
		if len(callArgNames) > 0 {
			b.WriteString(fmt.Sprintf("\tbaseArgs := []interface{}{%s}\n", strings.Join(callArgNames, ", ")))
		}

		// Count query.
		b.WriteString(fmt.Sprintf("\tcountSQL, countArgs := BuildCountQuery(%q, %q, %d, opts)\n", tableName, baseWhere, baseArgCount))
		if len(callArgNames) > 0 {
			b.WriteString("\tcountArgs = append(baseArgs, countArgs...)\n")
		}
		b.WriteString("\tvar total int64\n")
		b.WriteString("\tif err := m.conn().QueryRowContext(context.Background(), countSQL, countArgs...).Scan(&total); err != nil {\n")
		b.WriteString("\t\treturn nil, err\n")
		b.WriteString("\t}\n\n")

		// Select query.
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
		// Include loading — always applied (x-include is codegen metadata, not runtime option).
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

	case isSliceReturn:
		// Multi-row query without pagination: ([]Type, error)
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

	case isFind || seqType == "get":
		// Find method: (*Type, error)
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

	case seqType == "post":
		// Create method: (*Type, error)
		b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
		b.WriteString(fmt.Sprintf("\trow := m.conn().QueryRowContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
		b.WriteString(fmt.Sprintf("\treturn scan%s(row)\n", modelName))
		b.WriteString("}\n")

	case seqType == "put" || seqType == "delete":
		// Update/Delete: error
		b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
		b.WriteString(fmt.Sprintf("\t_, err := m.conn().ExecContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
		b.WriteString("\treturn err\n")
		b.WriteString("}\n")

	default:
		// Custom/unknown: determine from query cardinality or default to exec.
		if query != nil && query.Cardinality == "one" {
			b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
			b.WriteString(fmt.Sprintf("\trow := m.conn().QueryRowContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
			b.WriteString(fmt.Sprintf("\treturn scan%s(row)\n", modelName))
			b.WriteString("}\n")
		} else {
			b.WriteString(fmt.Sprintf("func (m *%s) %s(%s) %s {\n", implName, m.Name, m.ParamSig, m.ReturnSig))
			b.WriteString(fmt.Sprintf("\t_, err := m.conn().ExecContext(context.Background(),\n\t\t%q%s)\n", sqlStr, callArgs))
			b.WriteString("\treturn err\n")
			b.WriteString("}\n")
		}
	}
}
