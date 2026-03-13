package reporter

import (
	"fmt"
	"io"
)

// ChainLink mirrors orchestrator.ChainLink to avoid circular import.
type ChainLink struct {
	Kind      string
	File      string
	Line      int
	Summary   string
	Ownership string // "", "gen", "preserve"
}

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
		loc := link.File
		if link.Line > 0 {
			loc = fmt.Sprintf("%s:%d", link.File, link.Line)
		}
		fmt.Fprintf(w, "  %-10s %-45s %s\n", link.Kind, loc, link.Summary)
	}

	if len(artifactLinks) > 0 {
		fmt.Fprintf(w, "\n  ── Artifacts ──\n")
		for _, link := range artifactLinks {
			loc := link.File
			if link.Summary != "" && link.Summary != "(file)" {
				loc = link.File + ":" + link.Summary
			}
			ownerIcon := ""
			switch link.Ownership {
			case "preserve":
				ownerIcon = "preserve ✎"
			case "gen":
				ownerIcon = "gen"
			}
			fmt.Fprintf(w, "  %-10s %-45s %s\n", link.Kind, loc, ownerIcon)
		}
	}

	fmt.Fprintln(w)
}
