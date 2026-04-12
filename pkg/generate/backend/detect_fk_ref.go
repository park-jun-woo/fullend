//ff:func feature=rule type=util control=iteration dimension=1
//ff:what detectFKRef — 시퀀스가 이전 result 변수의 필드를 참조하는지 판정
package backend

import parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func detectFKRef(seq parsessac.Sequence, declared map[string]string) bool {
	if seq.Type != "get" {
		return false
	}
	for _, arg := range seq.Args {
		if arg.Source == "" || arg.Source == "request" || arg.Source == "currentUser" || arg.Source == "query" || arg.Source == "message" {
			continue
		}
		if _, ok := declared[arg.Source]; ok {
			return true
		}
	}
	return false
}
