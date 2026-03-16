//ff:func feature=ssac-gen type=generator control=iteration dimension=1
//ff:what 파생 인터페이스 배열을 Go 소스 코드(package model)로 렌더링
package generator

import "bytes"

func renderInterfaces(interfaces []derivedInterface) []byte {
	var buf bytes.Buffer
	buf.WriteString("package model\n\n")

	needTime := needsTimeImport(interfaces)
	needJSON := needsJSONImport(interfaces)
	needPagination := needsPaginationImport(interfaces)
	buf.WriteString("import (\n")
	buf.WriteString("\t\"database/sql\"\n")
	if needJSON {
		buf.WriteString("\t\"encoding/json\"\n")
	}
	if needTime {
		buf.WriteString("\t\"time\"\n")
	}
	if needPagination {
		buf.WriteString("\n\t\"github.com/geul-org/fullend/pkg/pagination\"\n")
	}
	buf.WriteString(")\n\n")

	for _, iface := range interfaces {
		renderSingleInterface(&buf, iface)
	}

	return buf.Bytes()
}
