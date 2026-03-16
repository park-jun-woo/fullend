//ff:func feature=ssac-gen type=generator control=sequence
//ff:what Go 파일 확장자 ".go"를 반환
package generator

// FileExtension은 Go 파일 확장자를 반환한다.
func (g *GoTarget) FileExtension() string { return ".go" }
