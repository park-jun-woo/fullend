//ff:func feature=rule type=loader control=sequence
//ff:what populateOpExtensions — x-sort, x-filter, x-pagination를 Ground에 등록
package ground

import (
	"encoding/json"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateOpExtensions(g *rule.Ground, opID string, ext map[string]any) {
	if raw, ok := ext["x-sort"]; ok {
		var xSort struct{ Allowed []string `json:"allowed"` }
		data, _ := json.Marshal(raw)
		if json.Unmarshal(data, &xSort) == nil {
			g.Lookup["OpenAPI.sort."+opID] = toSet(xSort.Allowed)
		}
	}
	if raw, ok := ext["x-filter"]; ok {
		var xFilter struct{ Allowed []string `json:"allowed"` }
		data, _ := json.Marshal(raw)
		if json.Unmarshal(data, &xFilter) == nil {
			g.Lookup["OpenAPI.filter."+opID] = toSet(xFilter.Allowed)
		}
	}
	if raw, ok := ext["x-pagination"]; ok {
		var xPag struct{ Style string `json:"style"` }
		data, _ := json.Marshal(raw)
		if json.Unmarshal(data, &xPag) == nil && xPag.Style != "" {
			g.Config["pagination."+opID] = true
		}
	}
}
