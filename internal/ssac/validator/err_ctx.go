//ff:type feature=ssac-validate type=util
//ff:what 검증 에러 컨텍스트 헬퍼
package validator

type errCtx struct {
	fileName string
	funcName string
	seqIndex int
}
