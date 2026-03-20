//ff:func feature=ssac-gen type=test control=sequence
//ff:what query arg가 있는 메서드에만 opts QueryOpts 파라미터를 추가하는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGenerateModelInterfaceQueryOptsExplicit(t *testing.T) {
	st := &validator.SymbolTable{
		Models: map[string]validator.ModelSymbol{
			"Reservation": {Methods: map[string]validator.MethodInfo{
				"ListByUserID": {Cardinality: "many"},
			}},
			"User": {Methods: map[string]validator.MethodInfo{
				"FindByID": {Cardinality: "one"},
			}},
		},
		DDLTables: map[string]validator.DDLTable{
			"reservations": {Columns: map[string]string{"id": "int64", "user_id": "int64"}},
			"users":        {Columns: map[string]string{"id": "int64"}},
		},
		Operations: map[string]validator.OperationSymbol{},
	}
	funcs := []parser.ServiceFunc{{
		Name: "ListMyReservations", FileName: "list_my_reservations.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "User.FindByID", Inputs: map[string]string{"ID": "currentUser.ID"}, Result: &parser.Result{Type: "User", Var: "user"}},
			{Type: parser.SeqGet, Model: "Reservation.ListByUserID", Inputs: map[string]string{"UserID": "currentUser.ID", "Opts": "query"}, Result: &parser.Result{Type: "[]Reservation", Var: "reservations"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"reservations": "reservations"}},
		},
	}}

	outDir := t.TempDir()
	if err := GenerateModelInterfaces(funcs, st, outDir); err != nil {
		t.Fatal(err)
	}

	data, err := readFile(t, outDir+"/model/models_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	// query arg가 있는 메서드에만 opts QueryOpts 추가
	assertContains(t, data, "ListByUserID(userID int64, opts QueryOpts)")
	// query arg가 없는 메서드에는 opts 없음
	assertNotContains(t, data, "FindByID(id int64, opts QueryOpts)")
	assertContains(t, data, "FindByID(id int64)")
}
