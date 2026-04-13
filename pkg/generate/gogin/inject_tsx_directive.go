//ff:func feature=gen-gogin type=generator control=sequence topic=output
//ff:what 단일 TSX 파일에 fullend:gen 디렉티브를 주입한다

package gogin

import (
	"os"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/contract"
)

// injectTSXDirective injects a fullend:gen directive into a single TSX file if absent.
func injectTSXDirective(path, entryName string) {
	src, err := os.ReadFile(path)
	if err != nil {
		return
	}
	content := string(src)

	if strings.Contains(content, "fullend:") {
		return
	}

	stmlName := strings.TrimSuffix(entryName, ".tsx") + ".html"
	ssotPath := "frontend/" + stmlName
	hash := contract.Hash7(content)

	d := &contract.Directive{Ownership: "gen", SSOT: ssotPath, Contract: hash}
	newContent := d.StringJS() + "\n" + content
	os.WriteFile(path, []byte(newContent), 0644)
}
