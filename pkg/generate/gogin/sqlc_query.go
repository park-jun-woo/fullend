//ff:type feature=gen-gogin type=model
//ff:what parsed sqlc query annotation

package gogin

// sqlcQuery represents a parsed sqlc query annotation.
type sqlcQuery struct {
	Name        string   // e.g. "FindByID"
	Cardinality string   // "one", "many", "exec"
	SQL         string   // the raw SQL string
	ParamCount  int      // number of $N placeholders
	Columns     []string // INSERT/UPDATE column names (for param mapping)
}
