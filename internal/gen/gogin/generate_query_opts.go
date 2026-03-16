//ff:func feature=gen-gogin type=generator control=sequence
//ff:what creates model/queryopts.go with parseQueryOpts, buildSelectQuery, buildCountQuery

package gogin

import (
	"os"
	"path/filepath"
)

// generateQueryOpts creates model/queryopts.go with parseQueryOpts, buildSelectQuery, buildCountQuery.
func generateQueryOpts(modelDir string) error {
	src := queryOptsTemplate()
	return os.WriteFile(filepath.Join(modelDir, "queryopts.go"), []byte(src), 0644)
}
