//ff:method feature=ssac-validate type=model
//ff:what ValidationError.Error() 메서드
package validator

import "fmt"

func (e ValidationError) Error() string {
	level := e.Level
	if level == "" {
		level = "ERROR"
	}
	return fmt.Sprintf("%s: %s:%s:seq[%d] %s — %s", level, e.FileName, e.FuncName, e.SeqIndex, e.Tag, e.Message)
}
