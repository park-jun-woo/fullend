//ff:type feature=gen-react type=model
//ff:what App.tsx 라우트 정보를 담는 구조체를 정의한다

package react

// route holds routing information for a single page component.
type route struct {
	path      string
	component string
	fileName  string
}
