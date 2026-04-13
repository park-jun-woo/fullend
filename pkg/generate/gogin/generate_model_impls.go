//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=interface-derive
//ff:what main dispatcher — orchestrates model implementation file generation

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

// generateModelImpls generates model implementation files that use database/sql directly.
func generateModelImpls(intDir string, models []string, modulePath, specsDir string, serviceFuncs []ssacparser.ServiceFunc, modelIncludeSpecs map[string][]string, cursorSpecs map[string]string) error {
	if len(models) == 0 {
		return nil
	}

	modelDir := filepath.Join(intDir, "model")
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return err
	}

	// Parse DDL files to get table/column info.
	tables := parseDDLFiles(specsDir)

	// Resolve per-model includes against DDL FK.
	includesByModel := make(map[string][]includeMapping)
	for modelName, specs := range modelIncludeSpecs {
		mappings, err := resolveIncludes(modelName, specs, tables)
		if err != nil {
			return fmt.Errorf("resolve includes for %s: %w", modelName, err)
		}
		if len(mappings) > 0 {
			includesByModel[modelName] = mappings
		}
	}

	// Parse query SQL files to get embedded SQL and metadata.
	queriesByModel := parseQueryFiles(specsDir)

	// Parse models_gen.go to get exact interface signatures.
	ifaceMethods := parseModelsGen(modelDir)

	// Collect per-model methods from service functions (for seq type info).
	seqTypeByModel := collectSeqTypes(serviceFuncs)

	// Generate types.go from DDL.
	if err := generateTypesFile(modelDir, models, tables, includesByModel); err != nil {
		return fmt.Errorf("types.go: %w", err)
	}

	// Generate queryopts.go (parseQueryOpts + SQL builders).
	if err := generateQueryOpts(modelDir); err != nil {
		return fmt.Errorf("queryopts.go: %w", err)
	}

	// Generate per-model implementation files.
	for _, m := range models {
		methods := ifaceMethods[m]
		table := tables[m]
		queries := queriesByModel[m]
		seqTypes := seqTypeByModel[m]
		if err := generateModelFile(modelDir, m, methods, table, queries, seqTypes, includesByModel[m], cursorSpecs); err != nil {
			return fmt.Errorf("%s.go: %w", strings.ToLower(m), err)
		}
	}

	// Generate include helpers if any model has includes.
	if len(includesByModel) > 0 {
		if err := generateIncludeHelpersFile(modelDir); err != nil {
			return fmt.Errorf("include_helpers.go: %w", err)
		}
	}

	return nil
}
