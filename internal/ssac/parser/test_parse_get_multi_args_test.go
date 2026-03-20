//ff:func feature=ssac-parse type=parser control=sequence
//ff:what TestParseGetMultiArgs: @get 다중 인자 파싱 검증
package parser

import "testing"

func TestParseGetMultiArgs(t *testing.T) {
	src := `package service

// @get []Reservation reservations = Reservation.ListByUserAndRoom({UserID: currentUser.ID, RoomID: request.RoomID})
func ListReservations(c *gin.Context) {}
`
	sfs := parseTestFile(t, src)
	seq := sfs[0].Sequences[0]
	assertEqual(t, "Result.Type", seq.Result.Type, "[]Reservation")
	if len(seq.Inputs) != 2 {
		t.Fatalf("expected 2 inputs, got %d", len(seq.Inputs))
	}
	assertEqual(t, "Inputs.UserID", seq.Inputs["UserID"], "currentUser.ID")
	assertEqual(t, "Inputs.RoomID", seq.Inputs["RoomID"], "request.RoomID")
}
