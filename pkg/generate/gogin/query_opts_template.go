//ff:type feature=gen-gogin type=model topic=output
//ff:what queryOptsTemplate — model/queryopts.go 의 정적 소스 (templates/query_opts.tmpl embed)

package gogin

import _ "embed"

//go:embed templates/query_opts.tmpl
var queryOptsTemplate string
