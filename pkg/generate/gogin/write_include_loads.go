//ff:func feature=gen-gogin type=generator control=iteration dimension=1
//ff:what writes include helper calls for loading related models

package gogin

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
)

// writeIncludeLoads writes include helper calls for a model.
func writeIncludeLoads(b *strings.Builder, includes []includeMapping) {
	for _, inc := range includes {
		helperName := "include" + strcase.ToGoPascal(inc.IncludeName)
		b.WriteString(fmt.Sprintf("\tif err := m.%s(items); err != nil {\n", helperName))
		b.WriteString("\t\treturn nil, err\n")
		b.WriteString("\t}\n")
	}
}
