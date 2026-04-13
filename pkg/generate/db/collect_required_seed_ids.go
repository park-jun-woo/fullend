//ff:func feature=gen-gogin type=util control=iteration dimension=3 topic=ddl
//ff:what collectRequiredSeedIDs — DEFAULT <int> + FK 감지해 필요 seed id 수집

package db

import (
	"fmt"
	"strconv"

	"github.com/park-jun-woo/fullend/pkg/parser/ddl"
)

// collectRequiredSeedIDs scans all columns for DEFAULT <int> REFERENCES <ref>(id).
// Key = "<ref>#<N>"; value always true.
func collectRequiredSeedIDs(tables []ddl.Table) map[string]bool {
	out := map[string]bool{}
	for _, t := range tables {
		fkByCol := map[string]ddl.ForeignKey{}
		for _, fk := range t.ForeignKeys {
			fkByCol[fk.Column] = fk
		}
		for col, def := range t.Defaults {
			n, err := strconv.ParseInt(def, 10, 64)
			if err != nil {
				continue
			}
			if fk, ok := fkByCol[col]; ok {
				out[fmt.Sprintf("%s#%d", fk.RefTable, n)] = true
			}
		}
	}
	return out
}
