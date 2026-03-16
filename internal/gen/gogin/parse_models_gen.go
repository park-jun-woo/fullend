//ff:func feature=gen-gogin type=parser control=iteration
//ff:what reads models_gen.go and extracts interface method signatures

package gogin

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// parseModelsGen reads models_gen.go and extracts interface method signatures.
// Returns map[ModelName][]ifaceMethod.
func parseModelsGen(modelDir string) map[string][]ifaceMethod {
	result := make(map[string][]ifaceMethod)

	path := filepath.Join(modelDir, "models_gen.go")
	data, err := os.ReadFile(path)
	if err != nil {
		return result
	}

	// Parse "type XxxModel interface {" blocks.
	ifaceRe := regexp.MustCompile(`type\s+(\w+)Model\s+interface\s*\{`)
	// Parse method lines: "MethodName(params) (returns)"
	methodRe := regexp.MustCompile(`^\s+(\w+)\(([^)]*)\)\s*(.+)$`)
	// Parse individual params: "name type"
	paramRe := regexp.MustCompile(`(\w+)\s+([\w.*\[\]]+)`)

	lines := strings.Split(string(data), "\n")
	var currentModel string

	for _, line := range lines {
		if m := ifaceRe.FindStringSubmatch(line); m != nil {
			currentModel = m[1]
			continue
		}
		if currentModel != "" && strings.TrimSpace(line) == "}" {
			currentModel = ""
			continue
		}
		if currentModel != "" {
			if m := methodRe.FindStringSubmatch(line); m != nil {
				method := ifaceMethod{
					Name:      m[1],
					ParamSig:  m[2],
					ReturnSig: m[3],
				}
				// Parse individual params.
				for _, pm := range paramRe.FindAllStringSubmatch(m[2], -1) {
					method.Params = append(method.Params, ifaceParam{Name: pm[1], Type: pm[2]})
				}
				result[currentModel] = append(result[currentModel], method)
			}
		}
	}

	return result
}
