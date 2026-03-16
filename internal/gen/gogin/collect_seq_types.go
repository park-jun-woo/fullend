//ff:func feature=gen-gogin type=util control=iteration
//ff:what extracts per-model method to sequence type mapping from service functions

package gogin

import (
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// collectSeqTypes extracts per-model method → sequence type mapping from service functions.
// Returns map[ModelName]map[MethodName]seqType.
func collectSeqTypes(funcs []ssacparser.ServiceFunc) map[string]map[string]string {
	result := make(map[string]map[string]string)

	for _, fn := range funcs {
		for _, seq := range fn.Sequences {
			// Skip @call — package-level funcs are not models.
			if seq.Type == "call" {
				continue
			}
			if seq.Model == "" {
				continue
			}
			parts := strings.SplitN(seq.Model, ".", 2)
			if len(parts) != 2 {
				continue
			}
			modelName := parts[0]
			methodName := parts[1]

			if result[modelName] == nil {
				result[modelName] = make(map[string]string)
			}
			result[modelName][methodName] = seq.Type
		}
	}

	return result
}
