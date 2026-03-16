//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=interface-derive
//ff:what 단일 인터페이스 정의를 버퍼에 출력
package generator

import (
	"bytes"
	"fmt"
)

func renderSingleInterface(buf *bytes.Buffer, iface derivedInterface) {
	fmt.Fprintf(buf, "type %s interface {\n", iface.Name)
	fmt.Fprintf(buf, "\tWithTx(tx *sql.Tx) %s\n", iface.Name)
	for _, m := range iface.Methods {
		params := renderMethodSignature(m)
		fmt.Fprintf(buf, "\t%s(%s) %s\n", m.Name, params, m.ReturnType)
	}
	buf.WriteString("}\n\n")
}
