//ff:func feature=gen-hurl type=util control=iteration
//ff:what Returns tables in deletion order (children before parents) via topological sort.
package hurl

import "sort"

// topoSortDelete returns tables in deletion order (children before parents).
func topoSortDelete(deps map[string]map[string]bool) []string {
	// Collect all tables.
	allTables := make(map[string]bool)
	for t := range deps {
		allTables[t] = true
		for dep := range deps[t] {
			allTables[dep] = true
		}
	}

	// Count how many tables depend on each table (in-degree as parent).
	childCount := make(map[string]int)
	for _, parents := range deps {
		for p := range parents {
			childCount[p]++
		}
	}

	// Start with leaf tables (no children depending on them).
	var queue []string
	for t := range allTables {
		if childCount[t] == 0 {
			queue = append(queue, t)
		}
	}
	sort.Strings(queue)

	var result []string
	visited := make(map[string]bool)

	for len(queue) > 0 {
		t := queue[0]
		queue = queue[1:]
		if visited[t] {
			continue
		}
		visited[t] = true
		result = append(result, t)

		// Decrement child count for parents of this table.
		for parent := range deps[t] {
			childCount[parent]--
			if childCount[parent] == 0 {
				queue = append(queue, parent)
				sort.Strings(queue)
			}
		}
	}

	// Add any remaining (circular deps — shouldn't happen).
	for t := range allTables {
		if !visited[t] {
			result = append(result, t)
		}
	}

	return result
}
