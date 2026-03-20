//ff:func feature=stml-validate type=test control=sequence
//ff:what 더미 스터디 프로젝트 전체 검증 통과 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestValidateDummyStudyPass(t *testing.T) {
	root := setupTestProject(t, dummyOpenAPI, nil, []string{"DatePicker"})
	pages := []parser.PageSpec{
		{Name: "login-page", FileName: "login-page.html", Actions: []parser.ActionBlock{{OperationID: "Login", Fields: []parser.FieldBind{{Name: "Email", Tag: "input", Type: "email"}, {Name: "Password", Tag: "input", Type: "password"}}}}},
		{Name: "my-reservations-page", FileName: "my-reservations-page.html", Fetches: []parser.FetchBlock{{OperationID: "ListMyReservations", Eaches: []parser.EachBlock{{Field: "reservations", Binds: []parser.FieldBind{{Name: "RoomID", Tag: "span"}}}}, States: []parser.StateBind{{Condition: "reservations.empty"}}}}, Actions: []parser.ActionBlock{{OperationID: "CreateReservation", Fields: []parser.FieldBind{{Name: "RoomID", Tag: "input", Type: "number"}, {Name: "StartAt", Tag: "data-component:DatePicker"}, {Name: "EndAt", Tag: "data-component:DatePicker"}}}}},
		{Name: "reservation-detail-page", FileName: "reservation-detail-page.html", Fetches: []parser.FetchBlock{{OperationID: "GetReservation", Params: []parser.ParamBind{{Name: "reservationId", Source: "route.ReservationID"}}, Binds: []parser.FieldBind{{Name: "reservation.Status", Tag: "span"}}}}},
	}
	errs := Validate(pages, root)
	if len(errs) > 0 { for _, e := range errs { t.Error(e.Error()) } }
}
