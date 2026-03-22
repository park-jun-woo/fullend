//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=model-collect
//ff:what x-include 스펙 하나를 파싱하고 FK 관계를 검증하여 includeMapping을 반환한다

package gogin

import (
	"fmt"
	"strings"
)

// resolveSingleInclude parses one x-include spec and validates it against DDL FK relationships.
func resolveSingleInclude(spec string, currentTable *ddlTable) (includeMapping, error) {
	colonIdx := strings.Index(spec, ":")
	if colonIdx <= 0 {
		return includeMapping{}, fmt.Errorf("invalid x-include format %q: expected 'column:table.column'", spec)
	}
	localColumn := spec[:colonIdx]
	targetRef := spec[colonIdx+1:]

	dotIdx := strings.Index(targetRef, ".")
	if dotIdx <= 0 {
		return includeMapping{}, fmt.Errorf("invalid x-include format %q: expected 'column:table.column'", spec)
	}
	targetTable := targetRef[:dotIdx]

	var fkCol *ddlColumn
	for i, col := range currentTable.Columns {
		if col.Name != localColumn {
			continue
		}
		if col.FKTable != targetTable {
			return includeMapping{}, fmt.Errorf("x-include %q: column %s.%s does not reference %s (references %q)",
				spec, currentTable.TableName, localColumn, targetTable, col.FKTable)
		}
		fkCol = &currentTable.Columns[i]
		break
	}
	if fkCol == nil {
		return includeMapping{}, fmt.Errorf("x-include %q: column %s not found in table %s",
			spec, localColumn, currentTable.TableName)
	}

	includeName := strings.TrimSuffix(localColumn, "_id")
	fieldName := fkColumnToFieldName(localColumn)
	targetModelName := singularize(targetTable)

	return includeMapping{
		IncludeName: includeName,
		FieldName:   fieldName,
		FieldType:   "*" + targetModelName,
		FKColumn:    localColumn,
		TargetTable: targetTable,
		TargetModel: targetModelName,
	}, nil
}
