//ff:func feature=stml-gen type=util control=iteration dimension=1 topic=import-collect
//ff:what FetchBlock에서 필요한 임포트(useParams, useState, 컴포넌트)를 수집한다
package generator

import (
	"strings"

	"github.com/park-jun-woo/fullend/internal/stml/parser"
)

func collectFetchImports(f parser.FetchBlock, is *importSet, compSet map[string]bool) {
	for _, p := range f.Params {
		if strings.HasPrefix(p.Source, "route.") {
			is.useParams = true
		}
	}
	for _, c := range f.Components {
		compSet[c.Name] = true
	}
	// Phase 5: infra params need useState
	if f.Paginate || f.Sort != nil || len(f.Filters) > 0 {
		is.useState = true
	}
	for _, child := range f.NestedFetches {
		collectFetchImports(child, is, compSet)
	}
}
