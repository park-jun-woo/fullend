//ff:func feature=crosscheck type=util control=sequence topic=config-check
//ff:what resolveClaimDDLColumn — JWT claim key 를 DDL 컬럼으로 매핑 (휴리스틱)

package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

// resolveClaimDDLColumn maps a JWT claim key to its likely DDL column Go type.
// Heuristic:
//  1. "<x>_id"  → try table "<x>s".id (e.g. user_id → users.id)
//  2. fallback: users.<claimKey> (e.g. org_id → users.org_id, email → users.email)
//  3. returns empty if no match
func resolveClaimDDLColumn(claimKey string, g *rule.Ground) (goType, tableColRef string, ok bool) {
	if strings.HasSuffix(claimKey, "_id") {
		plural := strings.TrimSuffix(claimKey, "_id") + "s"
		if gt, has := lookupTableColumn(plural, "id", g); has {
			return gt, plural + ".id", true
		}
	}
	if gt, has := lookupTableColumn("users", claimKey, g); has {
		return gt, "users." + claimKey, true
	}
	return "", "", false
}
