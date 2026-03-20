//ff:func feature=stml-parse type=test control=sequence
//ff:what Phase4 Tag + ClassName 파싱 검증
package parser

import ("strings"; "testing")

func TestPhase4_TagAndClassName(t *testing.T) {
	input := `<main class="max-w-4xl mx-auto p-6">
  <section data-fetch="ListMyReservations" class="mb-8">
    <ul data-each="reservations" class="space-y-3">
      <li class="flex justify-between p-4 border rounded">
        <span data-bind="RoomID" class="font-semibold"></span>
      </li>
    </ul>
  </section>
  <div data-action="CreateReservation" class="space-y-4">
    <input data-field="RoomID" type="number" placeholder="스터디룸 번호" class="w-full px-3 py-2 border rounded" />
    <button type="submit">예약하기</button>
  </div>
</main>`
	page, err := ParseReader("test.html", strings.NewReader(input))
	if err != nil { t.Fatal(err) }
	fetch := page.Fetches[0]
	if fetch.Tag != "section" { t.Errorf("Fetch.Tag = %q", fetch.Tag) }
	if fetch.ClassName != "mb-8" { t.Errorf("Fetch.ClassName = %q", fetch.ClassName) }
	each := fetch.Eaches[0]
	if each.Tag != "ul" { t.Errorf("Each.Tag = %q", each.Tag) }
	if each.ClassName != "space-y-3" { t.Errorf("Each.ClassName = %q", each.ClassName) }
	if each.ItemTag != "li" { t.Errorf("Each.ItemTag = %q", each.ItemTag) }
	if each.ItemClassName != "flex justify-between p-4 border rounded" { t.Errorf("Each.ItemClassName = %q", each.ItemClassName) }
	if each.Binds[0].ClassName != "font-semibold" { t.Errorf("Bind.ClassName = %q", each.Binds[0].ClassName) }
	action := page.Actions[0]
	if action.Tag != "div" { t.Errorf("Action.Tag = %q", action.Tag) }
	if action.ClassName != "space-y-4" { t.Errorf("Action.ClassName = %q", action.ClassName) }
	field := action.Fields[0]
	if field.Placeholder != "스터디룸 번호" { t.Errorf("Field.Placeholder = %q", field.Placeholder) }
	if field.ClassName != "w-full px-3 py-2 border rounded" { t.Errorf("Field.ClassName = %q", field.ClassName) }
	if action.SubmitText != "예약하기" { t.Errorf("SubmitText = %q", action.SubmitText) }
}
