//ff:func feature=gen-gogin type=util
//ff:what converts a plural table name to a singular model name

package gogin

import (
	"strings"

	"github.com/ettle/strcase"
	"github.com/jinzhu/inflection"
)

// singularize converts a plural table name to a singular model name.
// Rules: 'ies'-> 'y', 'sses'-> 'ss', 'xes'-> 'x', default strip trailing 's'.
func singularize(name string) string {
	singular := inflection.Singular(strings.ToLower(name))
	if len(singular) == 0 {
		return name
	}
	return strcase.ToGoPascal(singular)
}
