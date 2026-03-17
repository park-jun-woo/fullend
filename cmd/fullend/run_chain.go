//ff:func feature=cli type=command control=iteration dimension=1
//ff:what operationId 기반 feature chain 추적
package main

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/fullend/internal/orchestrator"
	"github.com/park-jun-woo/fullend/internal/reporter"
)

func runChain(operationID, specsDir string) {
	links, err := orchestrator.Chain(specsDir, operationID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	rlinks := make([]reporter.ChainLink, len(links))
	for i, l := range links {
		rlinks[i] = reporter.ChainLink{Kind: l.Kind, File: l.File, Line: l.Line, Summary: l.Summary, Ownership: l.Ownership}
	}
	reporter.PrintChain(os.Stdout, operationID, rlinks)
}
