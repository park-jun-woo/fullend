//ff:func feature=gen-gogin type=util control=iteration dimension=2
//ff:what reorders call arguments to match SQL column order from sqlcQuery

package gogin

// reorderCallArgs builds call arg names, reordering to match SQL column order.
// For INSERT/UPDATE, args follow column order; unmatched params are appended.
func reorderCallArgs(m ifaceMethod, query *sqlcQuery) []string {
	if query == nil || len(query.Columns) == 0 || len(m.Params) == 0 {
		var args []string
		for _, p := range m.Params {
			if p.Type != "QueryOpts" {
				args = append(args, p.Name)
			}
		}
		return args
	}

	paramByCol := make(map[string]string)
	for _, p := range m.Params {
		if p.Type == "QueryOpts" {
			continue
		}
		paramByCol[goToSnake(p.Name)] = p.Name
	}

	var callArgNames []string
	matched := make(map[string]bool)
	for _, col := range query.Columns {
		if paramName, ok := paramByCol[col]; ok {
			callArgNames = append(callArgNames, paramName)
			matched[paramName] = true
		}
	}
	for _, p := range m.Params {
		if p.Type == "QueryOpts" {
			continue
		}
		if !matched[p.Name] {
			callArgNames = append(callArgNames, p.Name)
		}
	}
	return callArgNames
}
