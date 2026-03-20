//ff:func feature=stml-gen type=test control=sequence
//ff:what Generate 결과의 dependencies 검증
package generator

import ("strings"; "testing"; "github.com/park-jun-woo/fullend/internal/stml/parser")

func TestGenerateResultDependencies(t *testing.T) {
	page, _ := parser.ParseReader("login-page.html", strings.NewReader(`<main>
  <div data-action="Login">
    <input data-field="Email" type="email" />
    <button type="submit">로그인</button>
  </div>
</main>`))
	outDir := t.TempDir()
	result, err := Generate([]parser.PageSpec{page}, "", outDir)
	if err != nil { t.Fatal(err) }
	if result.Pages != 1 { t.Errorf("expected 1 page, got %d", result.Pages) }
	if result.Dependencies["@tanstack/react-query"] != "^5" { t.Errorf("expected @tanstack/react-query ^5, got %q", result.Dependencies["@tanstack/react-query"]) }
	if result.Dependencies["react-hook-form"] != "^7" { t.Errorf("expected react-hook-form ^7, got %q", result.Dependencies["react-hook-form"]) }
}
