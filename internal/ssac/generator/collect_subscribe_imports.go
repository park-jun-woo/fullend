//ff:func feature=ssac-gen type=generator control=iteration dimension=1
//ff:what subscribe 함수에 필요한 import 경로를 수집
package generator

import (
	"sort"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// collectSubscribeImports는 subscribe 함수에 필요한 import를 수집한다.
func collectSubscribeImports(sf parser.ServiceFunc) []string {
	seen := map[string]bool{
		"context": true,
		"fmt":     true,
	}
	for _, seq := range sf.Sequences {
		if seq.Type == parser.SeqState {
			seen["states/"+seq.DiagramID+"state"] = true
		}
		if seq.Type == parser.SeqAuth {
			seen["authz"] = true
		}
		if seq.Type == parser.SeqPublish {
			seen["queue"] = true
		}
	}
	if needsCurrentUser(sf.Sequences) {
		seen["model"] = true
	}
	if hasWriteSequence(sf.Sequences) {
		seen["database/sql"] = true
	}
	for _, imp := range sf.Imports {
		seen[imp] = true
	}
	var imports []string
	order := []string{"context", "fmt"}
	for _, imp := range order {
		if seen[imp] {
			imports = append(imports, imp)
			delete(seen, imp)
		}
	}
	var dynamic []string
	for imp := range seen {
		dynamic = append(dynamic, imp)
	}
	sort.Strings(dynamic)
	return append(imports, dynamic...)
}
