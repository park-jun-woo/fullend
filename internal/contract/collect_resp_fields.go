//ff:func feature=contract type=util control=iteration dimension=1
//ff:what 시퀀스 목록에서 response 필드의 key:value 문자열을 수집한다
package contract

import (
	"sort"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// collectRespFields returns sorted "key:type" strings from response sequences.
func collectRespFields(seqs []ssacparser.Sequence) []string {
	var fields []string
	for _, seq := range seqs {
		if seq.Type != "response" || seq.Fields == nil {
			continue
		}
		keys := make([]string, 0, len(seq.Fields))
		for k := range seq.Fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fields = append(fields, k+":"+seq.Fields[k])
		}
	}
	return fields
}
