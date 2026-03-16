//ff:func feature=reporter type=formatter control=iteration dimension=1
//ff:what 기능 체인을 포맷팅하여 출력한다
package reporter

import (
	"fmt"
	"io"
)

// PrintChain writes a formatted feature chain to w.
func PrintChain(w io.Writer, operationID string, links []ChainLink) {
	fmt.Fprintf(w, "\n── Feature Chain: %s ──\n\n", operationID)

	if len(links) == 0 {
		fmt.Fprintln(w, "  No SSOT links found.")
		return
	}

	// Split SSOT links and artifact links.
	var ssotLinks, artifactLinks []ChainLink
	artifactKinds := map[string]bool{"Handler": true, "Model": true, "Authz": true, "States": true, "Types": true}
	for _, link := range links {
		if artifactKinds[link.Kind] {
			artifactLinks = append(artifactLinks, link)
		} else {
			ssotLinks = append(ssotLinks, link)
		}
	}

	for _, link := range ssotLinks {
		fmt.Fprintln(w, formatChainLink(link, false))
	}

	if len(artifactLinks) > 0 {
		fmt.Fprintf(w, "\n  ── Artifacts ──\n")
		for _, link := range artifactLinks {
			fmt.Fprintln(w, formatChainLink(link, true))
		}
	}

	fmt.Fprintln(w)
}
