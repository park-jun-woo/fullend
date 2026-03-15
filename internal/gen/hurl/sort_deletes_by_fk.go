//ff:func feature=gen-hurl type=util
//ff:what Delete ordering — FK dependency-aware topological sort for DELETE steps.
package hurl

import "sort"

// sortDeletesByFK sorts DELETE steps using DDL FK dependency graph.
// Children (tables with FK references) are deleted before parents.
func sortDeletesByFK(steps []scenarioStep, specsDir string) []scenarioStep {
	if len(steps) <= 1 || specsDir == "" {
		return steps
	}

	// Parse DDL to get FK relationships.
	tables := parseDDLFiles(specsDir)

	// Build dependency graph: table -> set of tables it depends on (via FK).
	// e.g. lessons depends on courses (lessons.course_id -> courses)
	deps := make(map[string]map[string]bool)
	for _, t := range tables {
		for _, fk := range t.FKTables {
			if deps[t.TableName] == nil {
				deps[t.TableName] = make(map[string]bool)
			}
			deps[t.TableName][fk] = true
		}
	}

	// Topological sort: tables with no dependents first, then their parents.
	// For DELETE order: children (FK holders) before parents.
	order := topoSortDelete(deps)

	// Map path -> table name for each DELETE step.
	tableOrder := make(map[string]int)
	for i, t := range order {
		tableOrder[t] = i
	}

	// Collect which tables have DELETE endpoints.
	deletableTables := make(map[string]bool)
	for _, s := range steps {
		t := inferTableFromPath(s.Path)
		if t != "" {
			deletableTables[t] = true
		}
	}

	// Build reverse deps: parent -> children that reference it.
	reverseDeps := make(map[string][]string)
	for child, parents := range deps {
		for parent := range parents {
			reverseDeps[parent] = append(reverseDeps[parent], child)
		}
	}

	// Filter out steps whose table has undeletable children (FK references with no DELETE endpoint).
	var filteredSteps []scenarioStep
	for _, s := range steps {
		t := inferTableFromPath(s.Path)
		if canDeleteTable(t, deletableTables, reverseDeps) {
			filteredSteps = append(filteredSteps, s)
		}
	}
	steps = filteredSteps

	// Sort steps by topological order.
	sort.SliceStable(steps, func(i, j int) bool {
		ti := inferTableFromPath(steps[i].Path)
		tj := inferTableFromPath(steps[j].Path)
		oi, oki := tableOrder[ti]
		oj, okj := tableOrder[tj]
		if oki && okj {
			return oi < oj
		}
		// Fallback: depth DESC, then path.
		if steps[i].PathDepth != steps[j].PathDepth {
			return steps[i].PathDepth > steps[j].PathDepth
		}
		return steps[i].Path < steps[j].Path
	})

	return steps
}
