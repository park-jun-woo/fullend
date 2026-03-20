//ff:func feature=stml-parse type=test control=sequence
//ff:what Phase4 State 텍스트 파싱 검증
package parser

import ("strings"; "testing")

func TestPhase4_StateText(t *testing.T) {
	input := `<main>
  <section data-fetch="ListMyReservations">
    <p data-state="reservations.empty" class="text-gray-400">예약이 없습니다</p>
  </section>
</main>`
	page, err := ParseReader("test.html", strings.NewReader(input))
	if err != nil { t.Fatal(err) }
	state := page.Fetches[0].States[0]
	if state.Tag != "p" { t.Errorf("State.Tag = %q", state.Tag) }
	if state.ClassName != "text-gray-400" { t.Errorf("State.ClassName = %q", state.ClassName) }
	if state.Text != "예약이 없습니다" { t.Errorf("State.Text = %q", state.Text) }
}
