//ff:func feature=ssac-gen type=generator control=sequence topic=subscribe
//ff:what subscribe 함수의 트랜잭션 커밋 코드를 출력
package ssac

import "bytes"

func writeSubTxCommit(buf *bytes.Buffer) {
	buf.WriteString("\tif err = tx.Commit(); err != nil {\n")
	buf.WriteString("\t\treturn fmt.Errorf(\"commit failed: %w\", err)\n")
	buf.WriteString("\t}\n\n")
}
