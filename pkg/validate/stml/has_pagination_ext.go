//ff:func feature=rule type=util control=sequence
//ff:what hasPaginationExt — operationId에 x-pagination 확장이 있는지 확인
package stml

import "github.com/park-jun-woo/fullend/pkg/rule"

func hasPaginationExt(ground *rule.Ground, opID string) bool {
	return ground.Config["pagination."+opID]
}
