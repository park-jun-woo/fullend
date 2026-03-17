//ff:func feature=funcspec type=parser control=iteration dimension=1
//ff:what FuncSpec의 빈 RequestFields/ResponseFields를 패키지 레벨 타입에서 보충한다
package funcspec

// fillMissingFields fills empty RequestFields/ResponseFields from
// companion struct files in the same directory.
func fillMissingFields(specs []FuncSpec, specDirs []string) {
	cache := make(map[string]map[string][]Field)
	for i := range specs {
		if len(specs[i].RequestFields) > 0 && len(specs[i].ResponseFields) > 0 {
			continue
		}
		dir := specDirs[i]
		typeMap, ok := cache[dir]
		if !ok {
			typeMap = collectPackageTypes(dir)
			cache[dir] = typeMap
		}
		fillSpecFromTypeMap(&specs[i], typeMap)
	}
}
