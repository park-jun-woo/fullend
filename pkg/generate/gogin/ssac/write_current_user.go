//ff:func feature=ssac-gen type=generator control=sequence topic=http-handler
//ff:what currentUser 추출 코드를 버퍼에 출력
package ssac

import "bytes"

func writeCurrentUser(buf *bytes.Buffer, needsCU bool) {
	if needsCU {
		var cuBuf bytes.Buffer
		goTemplates.ExecuteTemplate(&cuBuf, "currentUser", nil)
		buf.Write(cuBuf.Bytes())
		buf.WriteString("\n")
	}
}
