//ff:func feature=crosscheck type=util control=selection
//ff:what defaultErrStatus — 시퀀스 타입별 기본 에러 HTTP status 반환
package crosscheck

func defaultErrStatus(seqType string) int {
	switch seqType {
	case "empty":
		return 404
	case "exists":
		return 409
	case "state":
		return 409
	case "auth":
		return 403
	default:
		return 0
	}
}
