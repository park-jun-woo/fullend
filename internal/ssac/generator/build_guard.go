//ff:func feature=ssac-gen type=generator control=sequence topic=template-data
//ff:what empty/exists 가드의 제로값 비교 코드를 templateData에 설정
package generator

import "github.com/geul-org/fullend/internal/ssac/parser"

func buildGuard(d *templateData, seq parser.Sequence, resolver *FieldTypeResolver, resultTypes map[string]string) {
	d.Target = seq.Target
	if seq.Type != parser.SeqEmpty && seq.Type != parser.SeqExists {
		return
	}
	typeName := resolveGuardTypeName(seq.Target, resolver, resultTypes)
	d.ZeroCheck, d.ExistsCheck = zeroValueChecks(typeName)
}
