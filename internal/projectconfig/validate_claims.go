//ff:func feature=projectconfig type=util control=iteration dimension=1
//ff:what JWT 클레임 정의의 타입·예약어·중복을 검증한다
package projectconfig

import "fmt"

// validateClaims checks that all claim definitions use allowed types,
// do not use reserved JWT keys, and have no duplicate claim keys.
func validateClaims(auth *Auth) error {
	usedKeys := make(map[string]string) // claim_key → field_name (for duplicate detection)
	for field, def := range auth.Claims {
		if !allowedClaimTypes[def.GoType] {
			return fmt.Errorf("fullend.yaml: auth.claims.%s — type %q is not allowed (allowed: string, int64, bool)", field, def.GoType)
		}
		if jwtReservedKeys[def.Key] {
			return fmt.Errorf("fullend.yaml: auth.claims.%s — claim key %q is a reserved JWT key", field, def.Key)
		}
		if prev, dup := usedKeys[def.Key]; dup {
			return fmt.Errorf("fullend.yaml: auth.claims — duplicate claim key %q (used by %s and %s)", def.Key, prev, field)
		}
		usedKeys[def.Key] = field
	}
	return nil
}
