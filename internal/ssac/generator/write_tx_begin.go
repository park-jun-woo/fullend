//ff:func feature=ssac-gen type=generator control=sequence topic=http-handler
//ff:what 트랜잭션 시작 코드를 버퍼에 출력
package generator

import "bytes"

func writeTxBegin(buf *bytes.Buffer) {
	buf.WriteString("\ttx, err := h.DB.BeginTx(c.Request.Context(), nil)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\tc.JSON(http.StatusInternalServerError, gin.H{\"error\": \"transaction failed\"})\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tdefer tx.Rollback()\n\n")
}
