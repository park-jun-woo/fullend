//ff:func feature=stml-parse type=test control=sequence
//ff:what 스터디룸 수정 페이지 HTML 파싱 검증
package parser

import ("strings"; "testing")

func TestParseRoomEditPage(t *testing.T) {
	input := `<main>
  <div data-action="UpdateRoom" data-param-room-id="route.RoomID">
    <input data-field="Name" />
    <input data-field="Capacity" type="number" />
    <input data-field="Location" />
    <button type="submit">수정</button>
  </div>
  <footer data-state="canDelete">
    <button data-action="DeleteRoom" data-param-room-id="route.RoomID">스터디룸 삭제</button>
  </footer>
</main>`
	page, err := ParseReader("room-edit-page.html", strings.NewReader(input))
	if err != nil { t.Fatal(err) }
	if len(page.Actions) != 2 { t.Fatalf("Actions = %d, want 2", len(page.Actions)) }
	update := page.Actions[0]
	if update.OperationID != "UpdateRoom" { t.Errorf("OperationID = %q", update.OperationID) }
	if len(update.Params) != 1 { t.Fatalf("Params = %d", len(update.Params)) }
	if update.Params[0].Name != "roomId" { t.Errorf("Param.Name = %q", update.Params[0].Name) }
	if len(update.Fields) != 3 { t.Fatalf("Fields = %d", len(update.Fields)) }
	del := page.Actions[1]
	if del.OperationID != "DeleteRoom" { t.Errorf("OperationID = %q", del.OperationID) }
	if len(del.Params) != 1 { t.Fatalf("Params = %d", len(del.Params)) }
}
