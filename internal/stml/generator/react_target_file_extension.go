//ff:func feature=stml-gen type=generator control=sequence
//ff:what ReactTarget의 파일 확장자 .tsx를 반환한다
package generator

func (r *ReactTarget) FileExtension() string {
	return ".tsx"
}
