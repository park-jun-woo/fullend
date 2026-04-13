//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=ddl
//ff:what buildAutoNobodySeeds — DEFAULT N FK 패턴 감지해 nobody seed INSERT 자동 생성

package db

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/ddl"
)

// buildAutoNobodySeeds emits INSERT seeds for required <refTable, N> pairs.
// Values satisfy CHECK enums + NOT NULL; ON CONFLICT DO NOTHING for idempotency.
func buildAutoNobodySeeds(tables []ddl.Table, tableByName map[string]*ddl.Table) string {
	required := collectRequiredSeedIDs(tables)
	if len(required) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, key := range sortedSeedKeys(required) {
		ref, idVal := parseSeedKey(key)
		t := tableByName[ref]
		if t == nil || seedAlreadyExists(t, idVal) {
			continue
		}
		sb.WriteString(buildSeedInsertStmt(t, idVal))
		sb.WriteString("\n")
	}
	return sb.String()
}
