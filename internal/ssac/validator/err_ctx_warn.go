//ff:method feature=ssac-validate type=util
//ff:what errCtx에서 WARNING 레벨 ValidationError를 생성
package validator

func (c errCtx) warn(tag, msg string) ValidationError {
	return ValidationError{
		FileName: c.fileName,
		FuncName: c.funcName,
		SeqIndex: c.seqIndex,
		Tag:      tag,
		Message:  msg,
		Level:    "WARNING",
	}
}
