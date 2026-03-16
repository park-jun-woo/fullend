//ff:func feature=ssac-gen type=generator control=sequence
//ff:what 트랜잭션 커밋 코드를 버퍼에 출력
package generator

import "bytes"

func writeTxCommit(buf *bytes.Buffer) {
	buf.WriteString("\tif err = tx.Commit(); err != nil {\n")
	buf.WriteString("\t\tc.JSON(http.StatusInternalServerError, gin.H{\"error\": \"commit failed\"})\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n\n")
}
