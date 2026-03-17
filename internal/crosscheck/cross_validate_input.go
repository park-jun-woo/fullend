//ff:type feature=crosscheck type=model
//ff:what 교차 검증에 필요한 파싱된 SSOT 입력 데이터 구조체
package crosscheck

import (
	"github.com/park-jun-woo/fullend/internal/genapi"
	"github.com/park-jun-woo/fullend/internal/projectconfig"
)

// CrossValidateInput holds the pre-loaded data from individual validations.
type CrossValidateInput struct {
	*genapi.ParsedSSOTs
	DTOTypes        map[string]bool                    // model types marked with @dto (skip DDL matching)
	Middleware      []string                           // from fullend.yaml backend.middleware
	Archived        *ArchivedInfo                      // @archived tables/columns from DDL
	Claims          map[string]projectconfig.ClaimDef  // from fullend.yaml backend.auth.claims
	QueueBackend    string                             // from fullend.yaml queue.backend ("postgres", "memory", "")
	AuthzPackage    string                             // from fullend.yaml authz.package ("" = default pkg/authz)
	SensitiveCols   map[string]map[string]bool         // @sensitive columns per table (table → column → true)
	NoSensitiveCols map[string]map[string]bool         // @nosensitive columns per table (suppress WARNING)
	Roles           []string                           // from fullend.yaml auth.roles
}
