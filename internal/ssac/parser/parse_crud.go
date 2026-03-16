//ff:func feature=ssac-parse type=parser control=sequence
//ff:what @get/@post/@put/@delete CRUD 시퀀스 파싱
package parser

import "strings"

// parseCRUD는 @get/@post/@put/@delete를 파싱한다.
// hasResult=true: Type var = Model.Method({Key: val, ...})
// hasResult=false: Model.Method({Key: val, ...})
func parseCRUD(seqType, rest string, hasResult bool) (*Sequence, error) {
	rest = strings.TrimSpace(rest)
	seq := &Sequence{Type: seqType}

	if hasResult {
		// Type var = Model.Method({Key: val, ...})
		eqIdx := strings.Index(rest, "=")
		if eqIdx < 0 {
			return nil, nil
		}
		lhs := strings.TrimSpace(rest[:eqIdx])
		rhs := strings.TrimSpace(rest[eqIdx+1:])

		result := parseResult(lhs)
		if result == nil {
			return nil, nil
		}
		seq.Result = result

		model, inputs, _, err := parseCallExprInputs(rhs)
		if err != nil {
			return nil, err
		}
		seq.Package, seq.Model = splitPackagePrefix(model)
		seq.Inputs = inputs
	} else {
		// Model.Method({Key: val, ...})
		model, inputs, _, err := parseCallExprInputs(rest)
		if err != nil {
			return nil, err
		}
		seq.Package, seq.Model = splitPackagePrefix(model)
		seq.Inputs = inputs
	}

	return seq, nil
}
