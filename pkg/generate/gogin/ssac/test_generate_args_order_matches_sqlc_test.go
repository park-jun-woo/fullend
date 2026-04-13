//ff:func feature=ssac-gen type=test control=sequence
//ff:what sqlc 파라미터 순서대로 인자를 생성하는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGenerateArgsOrderMatchesSqlc(t *testing.T) {
	// sqlc: UPDATE gigs SET status = $1 WHERE id = $2
	// → Params: ["Status", "ID"] (SQL 순서)
	// SSaC: @put Gig.UpdateStatus({ID: request.GigID, Status: "published"})
	// 알파벳순이면: gigModel.UpdateStatus(gigID, "published") — ID < Status (잘못됨)
	// SQL 순서:    gigModel.UpdateStatus("published", gigID) — $1=status, $2=id (올바름)
	st := &validator.SymbolTable{
		Models: map[string]validator.ModelSymbol{
			"Gig": {Methods: map[string]validator.MethodInfo{
				"UpdateStatus": {Cardinality: "exec", Params: []string{"Status", "ID"}},
			}},
		},
		Operations: map[string]validator.OperationSymbol{},
		DDLTables:  map[string]validator.DDLTable{},
	}
	sf := ssacparser.ServiceFunc{
		Name: "PublishGig", FileName: "publish_gig.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqPut, Model: "Gig.UpdateStatus", Inputs: map[string]string{"ID": "request.GigID", "Status": `"published"`}},
		},
	}
	code := mustGenerate(t, sf, st)
	// SQL 순서: $1=status, $2=id → ("published", gigID)
	assertContains(t, code, `h.GigModel.WithTx(tx).UpdateStatus("published", gigID)`)
}
