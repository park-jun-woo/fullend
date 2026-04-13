//ff:func feature=ssac-gen type=generator control=sequence topic=subscribe
//ff:what subscribe 함수의 트랜잭션 시작 코드를 출력
package ssac

import "bytes"

func writeSubTxBegin(buf *bytes.Buffer) {
	buf.WriteString("\ttx, err := h.DB.BeginTx(ctx, nil)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn fmt.Errorf(\"transaction failed: %w\", err)\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tdefer tx.Rollback()\n\n")
}
