//ff:func feature=ssac-gen type=generator control=sequence
//ff:what CRUD 시퀀스의 모델 호출 코드를 templateData에 설정
package generator

func buildCRUDModelData(d *templateData, parts []string, useTx bool) {
	modelRef := "h." + toGoPascal(parts[0]) + "Model"
	if useTx {
		modelRef += ".WithTx(tx)"
	}
	d.ModelCall = modelRef + "." + parts[1]
}
