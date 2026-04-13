//ff:func feature=ssac-gen type=generator control=sequence topic=template-data
//ff:what empty/exists 가드의 제로값 비교 코드를 templateData에 설정
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func buildGuard(d *templateData, seq ssacparser.Sequence, resolver *FieldTypeResolver, resultTypes map[string]string) {
	d.Target = seq.Target
	if seq.Type != ssacparser.SeqEmpty && seq.Type != ssacparser.SeqExists {
		return
	}
	typeName := resolveGuardTypeName(seq.Target, resolver, resultTypes)
	d.ZeroCheck, d.ExistsCheck = zeroValueChecks(typeName)
}
