//ff:func feature=stml-parse type=parser control=sequence
//ff:what io.Reader에서 HTML을 파싱하여 PageSpec 반환
package parser

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// ParseReader parses HTML from a reader and returns a PageSpec.
func ParseReader(filename string, r io.Reader) (PageSpec, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return PageSpec{}, fmt.Errorf("html parse: %w", err)
	}

	name := strings.TrimSuffix(filename, ".html")
	page := PageSpec{
		Name:     name,
		FileName: filename,
	}

	walkTopLevel(doc, &page)
	return page, nil
}
