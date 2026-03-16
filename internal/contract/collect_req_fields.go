//ff:func feature=contract type=util control=iteration dimension=1
//ff:what 시퀀스 목록에서 request 소스 필드명을 수집한다
package contract

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

// collectReqFields returns request field names from all sequences.
func collectReqFields(seqs []ssacparser.Sequence) []string {
	var fields []string
	for _, seq := range seqs {
		fields = append(fields, collectReqFieldsFromSeq(seq.Args)...)
	}
	return fields
}
