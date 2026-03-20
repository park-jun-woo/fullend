//ff:func feature=stml-parse type=test control=sequence
//ff:what 예약 상세 페이지 HTML 파싱 검증
package parser

import ("strings"; "testing")

func TestParseReservationDetailPage(t *testing.T) {
	input := `<main>
  <article data-fetch="GetReservation" data-param-reservation-id="route.ReservationID">
    <span data-bind="reservation.Status"></span>
    <dd data-bind="reservation.RoomID"></dd>
    <dd data-bind="reservation.StartAt"></dd>
    <dd data-bind="reservation.EndAt"></dd>
    <footer data-state="canCancel">
      <button data-action="CancelReservation" data-param-reservation-id="route.ReservationID">예약 취소</button>
    </footer>
  </article>
</main>`
	page, err := ParseReader("reservation-detail-page.html", strings.NewReader(input))
	if err != nil { t.Fatal(err) }
	if len(page.Fetches) != 1 { t.Fatalf("Fetches = %d, want 1", len(page.Fetches)) }
	fetch := page.Fetches[0]
	if fetch.OperationID != "GetReservation" { t.Errorf("OperationID = %q", fetch.OperationID) }
	if len(fetch.Params) != 1 { t.Fatalf("Params = %d, want 1", len(fetch.Params)) }
	if fetch.Params[0].Name != "reservationId" { t.Errorf("Param.Name = %q", fetch.Params[0].Name) }
	if fetch.Params[0].Source != "route.ReservationID" { t.Errorf("Param.Source = %q", fetch.Params[0].Source) }
	if len(fetch.Binds) != 4 { t.Errorf("Binds = %d, want 4", len(fetch.Binds)) }
	if len(fetch.States) != 1 { t.Fatalf("States = %d, want 1", len(fetch.States)) }
	if fetch.States[0].Condition != "canCancel" { t.Errorf("State = %q", fetch.States[0].Condition) }
}
