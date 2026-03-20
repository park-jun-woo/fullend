//ff:func feature=stml-parse type=test control=sequence
//ff:what Phase4 State 내 Action 자식 파싱 검증
package parser

import ("strings"; "testing")

func TestPhase4_StateWithAction(t *testing.T) {
	input := `<main>
  <article data-fetch="GetReservation" data-param-reservation-id="route.ReservationID">
    <footer data-state="canCancel" class="mt-8 pt-4 border-t">
      <button data-action="CancelReservation" data-param-reservation-id="route.ReservationID">예약 취소</button>
    </footer>
  </article>
</main>`
	page, err := ParseReader("test.html", strings.NewReader(input))
	if err != nil { t.Fatal(err) }
	state := page.Fetches[0].States[0]
	if state.Tag != "footer" { t.Errorf("State.Tag = %q", state.Tag) }
	if state.ClassName != "mt-8 pt-4 border-t" { t.Errorf("State.ClassName = %q", state.ClassName) }
	if len(state.Children) != 1 { t.Fatalf("State.Children = %d", len(state.Children)) }
	if state.Children[0].Kind != "action" { t.Errorf("Kind = %q", state.Children[0].Kind) }
	if state.Children[0].Action.OperationID != "CancelReservation" { t.Errorf("OperationID = %q", state.Children[0].Action.OperationID) }
}
