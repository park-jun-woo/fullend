//ff:func feature=ssac-parse type=parser control=sequence
//ff:what TestParseStructs: SSaC 파일 내 struct 정의 파싱 후 이름·필드 검증
package parser

import "testing"

func TestParseStructs(t *testing.T) {
	src := `package service

type OnOrderCompletedMessage struct {
	OrderID int64
	Email   string
	Amount  int64
}

// @subscribe "order.completed"
// @get Order order = Order.FindByID({ID: message.OrderID})
func OnOrderCompleted(message OnOrderCompletedMessage) {}
`
	sfs := parseTestFile(t, src)
	sf := sfs[0]
	if len(sf.Structs) != 1 {
		t.Fatalf("expected 1 struct, got %d", len(sf.Structs))
	}
	si := sf.Structs[0]
	assertEqual(t, "Struct.Name", si.Name, "OnOrderCompletedMessage")
	if len(si.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(si.Fields))
	}
	assertEqual(t, "Field0.Name", si.Fields[0].Name, "OrderID")
	assertEqual(t, "Field0.Type", si.Fields[0].Type, "int64")
	assertEqual(t, "Field1.Name", si.Fields[1].Name, "Email")
	assertEqual(t, "Field1.Type", si.Fields[1].Type, "string")
}
