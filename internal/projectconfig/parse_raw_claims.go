//ff:func feature=projectconfig type=util control=iteration dimension=1
//ff:what RawClaims 맵을 ClaimDef 맵으로 변환한다
package projectconfig

import "strings"

// parseRawClaims converts RawClaims (field → "claim_key" or "claim_key:go_type")
// into parsed ClaimDef map.
func parseRawClaims(rawClaims map[string]string) map[string]ClaimDef {
	claims := make(map[string]ClaimDef, len(rawClaims))
	for field, raw := range rawClaims {
		parts := strings.SplitN(raw, ":", 2)
		def := ClaimDef{Key: parts[0], GoType: "string"}
		if len(parts) == 2 && parts[1] != "" {
			def.GoType = parts[1]
		}
		claims[field] = def
	}
	return claims
}
