//ff:func feature=ssac-gen type=generator control=sequence
//ff:what 패키지, import, 본문을 조립하여 gofmt 적용한 Go 소스를 반환
package generator

import (
	"bytes"
	"fmt"
	"go/format"
)

func assembleGoSource(pkgName string, imports []string, body []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("package " + pkgName + "\n\n")
	if len(imports) > 0 {
		buf.WriteString("import (\n")
		for _, imp := range imports {
			fmt.Fprintf(&buf, "\t%q\n", imp)
		}
		buf.WriteString(")\n\n")
	}
	buf.Write(body)

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.Bytes(), fmt.Errorf("gofmt 실패: %w\n--- raw ---\n%s", err, buf.String())
	}
	return formatted, nil
}
