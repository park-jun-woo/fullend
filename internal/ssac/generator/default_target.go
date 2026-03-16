//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what Go Target 기본 인스턴스를 반환
package generator

// DefaultTarget은 Go Target을 반환한다.
func DefaultTarget() Target {
	return &GoTarget{}
}
