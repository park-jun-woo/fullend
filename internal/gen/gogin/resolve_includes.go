//ff:func feature=gen-gogin type=util control=iteration dimension=2 topic=model-collect
//ff:what resolves x-include specs against DDL FK relationships

package gogin

// resolveIncludes resolves x-include specs against DDL FK relationships.
// Format: "column:table.column" (e.g. "instructor_id:users.id"). Forward FK only.
func resolveIncludes(modelName string, includeSpecs []string, tables map[string]*ddlTable) ([]includeMapping, error) {
	currentTable := tables[modelName]
	if currentTable == nil {
		return nil, nil
	}

	var mappings []includeMapping

	for _, spec := range includeSpecs {
		m, err := resolveSingleInclude(spec, currentTable)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, m)
	}

	return mappings, nil
}
