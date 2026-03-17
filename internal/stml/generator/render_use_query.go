//ff:func feature=stml-gen type=generator control=sequence
//ff:what FetchBlock에 대한 useQuery 훅 호출 코드를 생성한다
package generator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/stml/parser"
)

// renderUseQuery generates a useQuery hook call.
func renderUseQuery(f parser.FetchBlock) string {
	alias := toLowerFirst(f.OperationID) + "Data"
	paramValues := renderParamValues(f.Params)
	paramArgs := renderParamArgs(f.Params)

	// queryKey parts
	queryKey := fmt.Sprintf("'%s'", f.OperationID)
	if paramValues != "" {
		queryKey += ", " + paramValues
	}
	// Phase 5: infra params in queryKey
	if f.Paginate {
		queryKey += ", page, limit"
	}
	if f.Sort != nil {
		queryKey += ", sortBy, sortDir"
	}
	if len(f.Filters) > 0 {
		queryKey += ", filters"
	}

	// API call args
	hasInfra := f.Paginate || f.Sort != nil || len(f.Filters) > 0
	apiArgs := paramArgs
	if hasInfra {
		apiArgs = renderInfraApiArgs(f, paramArgs)
	}
	apiCall := fmt.Sprintf("api.%s(%s)", f.OperationID, apiArgs)

	return fmt.Sprintf(`const { data: %s, isLoading: %sLoading, error: %sError } = useQuery({
    queryKey: [%s],
    queryFn: () => %s,
  })`, alias, alias, alias, queryKey, apiCall)
}
