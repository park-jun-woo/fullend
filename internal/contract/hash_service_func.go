//ff:func feature=contract type=util control=sequence
//ff:what SSaC 서비스 함수의 계약 해시를 계산한다
package contract

import (
	"sort"
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// HashServiceFunc computes a contract hash for an SSaC service function.
// The hash is derived from: operationId + sequence types + request fields + response fields.
func HashServiceFunc(sf ssacparser.ServiceFunc) string {
	var parts []string
	parts = append(parts, sf.Name)

	// sequence types in order
	parts = append(parts, strings.Join(collectSeqTypes(sf.Sequences), ","))

	// request args (fields from request source)
	reqFields := collectReqFields(sf.Sequences)
	sort.Strings(reqFields)
	parts = append(parts, strings.Join(reqFields, ","))

	// response fields
	parts = append(parts, strings.Join(collectRespFields(sf.Sequences), ","))

	return Hash7(strings.Join(parts, "|"))
}
