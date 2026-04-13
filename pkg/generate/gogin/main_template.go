//ff:func feature=gen-gogin type=generator control=sequence topic=output
//ff:what mainWithDomainsTemplate — text/template + embed 로 cmd/main.go 렌더

package gogin

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed templates/main.tmpl
var mainTmplSrc string

var mainTmpl = template.Must(template.New("main").Option("missingkey=zero").Parse(mainTmplSrc))

// mainWithDomainsTemplate renders the domain-mode cmd/main.go source.
// Placeholder fields may be empty when a feature (auth/queue/...) is disabled.
func mainWithDomainsTemplate(osImport, importBlock, queueImport, builtinImport, jwtFlagLine, authzBlock, queueInitBlock, builtinInitBlock, initBlock, queueSubscribeBlock, dbName string) string {
	data := MainTmplData{
		OsImport:            osImport,
		ImportBlock:         importBlock,
		QueueImport:         queueImport,
		BuiltinImport:       builtinImport,
		JWTFlagLine:         jwtFlagLine,
		AuthzBlock:          authzBlock,
		QueueInitBlock:      queueInitBlock,
		BuiltinInitBlock:    builtinInitBlock,
		InitBlock:           initBlock,
		QueueSubscribeBlock: queueSubscribeBlock,
		DBName:              dbName,
	}
	var buf bytes.Buffer
	if err := mainTmpl.Execute(&buf, data); err != nil {
		return "// template execute error: " + err.Error()
	}
	return buf.String()
}
