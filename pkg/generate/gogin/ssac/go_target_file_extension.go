//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what Go 파일 확장자 ".go"를 반환
package ssac

// FileExtension은 Go 파일 확장자를 반환한다.
func (g *GoTarget) FileExtension() string { return ".go" }
