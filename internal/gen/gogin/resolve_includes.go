//ff:func feature=gen-gogin type=util control=iteration dimension=2
//ff:what resolves x-include specs against DDL FK relationships

package gogin

import (
	"fmt"
	"strings"
)

// resolveIncludes resolves x-include specs against DDL FK relationships.
// Format: "column:table.column" (e.g. "instructor_id:users.id"). Forward FK only.
func resolveIncludes(modelName string, includeSpecs []string, tables map[string]*ddlTable) ([]includeMapping, error) {
	currentTable := tables[modelName]
	if currentTable == nil {
		return nil, nil
	}

	var mappings []includeMapping

	for _, spec := range includeSpecs {
		// Parse "instructor_id:users.id"
		colonIdx := strings.Index(spec, ":")
		if colonIdx <= 0 {
			return nil, fmt.Errorf("invalid x-include format %q: expected 'column:table.column'", spec)
		}
		localColumn := spec[:colonIdx]
		targetRef := spec[colonIdx+1:]

		dotIdx := strings.Index(targetRef, ".")
		if dotIdx <= 0 {
			return nil, fmt.Errorf("invalid x-include format %q: expected 'column:table.column'", spec)
		}
		targetTable := targetRef[:dotIdx]

		// Validate: localColumn exists in current table with FK to targetTable.
		var fkCol *ddlColumn
		for i, col := range currentTable.Columns {
			if col.Name != localColumn {
				continue
			}
			if col.FKTable != targetTable {
				return nil, fmt.Errorf("x-include %q: column %s.%s does not reference %s (references %q)",
					spec, currentTable.TableName, localColumn, targetTable, col.FKTable)
			}
			fkCol = &currentTable.Columns[i]
			break
		}
		if fkCol == nil {
			return nil, fmt.Errorf("x-include %q: column %s not found in table %s",
				spec, localColumn, currentTable.TableName)
		}

		includeName := strings.TrimSuffix(localColumn, "_id")
		fieldName := fkColumnToFieldName(localColumn)
		targetModelName := singularize(targetTable)

		mappings = append(mappings, includeMapping{
			IncludeName: includeName,
			FieldName:   fieldName,
			FieldType:   "*" + targetModelName,
			FKColumn:    localColumn,
			TargetTable: targetTable,
			TargetModel: targetModelName,
		})
	}

	return mappings, nil
}
