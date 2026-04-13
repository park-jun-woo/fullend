//ff:func feature=gen-react type=util control=sequence
//ff:what tanstack 유무에 따라 main.tsx 소스 문자열을 반환한다

package react

// mainTSXSource returns the main.tsx source string based on whether tanstack is used.
func mainTSXSource(useTanstack bool) string {
	if useTanstack {
		return mainTSXWithTanstack
	}
	return mainTSXWithoutTanstack
}
