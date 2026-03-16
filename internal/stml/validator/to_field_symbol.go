//ff:func feature=stml-validate type=util control=sequence
//ff:what OpenAPI 스키마를 FieldSymbol로 변환
package validator

func toFieldSymbol(s openAPISchema, schemas map[string]openAPISchema) FieldSymbol {
	if s.Ref != "" {
		return FieldSymbol{Type: "object", ItemType: refName(s.Ref)}
	}
	if s.Type == "array" && s.Items != nil {
		itemType := s.Items.Type
		if s.Items.Ref != "" {
			itemType = refName(s.Items.Ref)
		}
		return FieldSymbol{Type: "array", ItemType: itemType}
	}
	return FieldSymbol{Type: s.Type}
}
