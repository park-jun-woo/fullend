//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=policy-check
//ff:what @ownership 어노테이션의 테이블·컬럼이 DDL에 존재하는지 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/internal/policy"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func checkOwnershipDDL(allOwnerships []policy.OwnershipMapping, st *ssacvalidator.SymbolTable) []CrossError {
	var errs []CrossError
	for _, om := range allOwnerships {
		errs = append(errs, checkSingleOwnershipDDL(om, st)...)
	}
	return errs
}
