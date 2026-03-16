//ff:func feature=gen-gogin type=util control=iteration dimension=1
//ff:what builds a map from generated .go filename to operationID

package gogin

import (
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// buildFileToOperationID builds a map from generated .go filename to operationID.
// SSaC parser sets FileName to the .ssac basename (e.g. "create_gig.ssac"),
// so we convert to .go extension for lookup by the generated file.
func buildFileToOperationID(funcs []ssacparser.ServiceFunc) map[string]string {
	result := make(map[string]string)
	for _, fn := range funcs {
		key := strings.TrimSuffix(fn.FileName, ".ssac") + ".go"
		result[key] = fn.Name
	}
	return result
}
