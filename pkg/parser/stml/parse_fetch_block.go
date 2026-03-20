//ff:func feature=stml-parse type=parser control=iteration dimension=1
//ff:what data-fetch 요소에서 FetchBlock 구성
package parser

import "golang.org/x/net/html"

// parseFetchBlock builds a FetchBlock from a data-fetch element and its descendants.
func parseFetchBlock(n *html.Node, operationID string) FetchBlock {
	fb := FetchBlock{
		Tag:         n.Data,
		ClassName:   getAttr(n, "class"),
		OperationID: operationID,
		Params:      extractParams(n),
	}

	if hasAttr(n, "data-paginate") {
		fb.Paginate = true
	}
	if v := getAttr(n, "data-sort"); v != "" {
		fb.Sort = parseSortDecl(v)
	}
	if v := getAttr(n, "data-filter"); v != "" {
		fb.Filters = splitTrim(v)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkFetchChildren(c, &fb)
	}
	return fb
}
