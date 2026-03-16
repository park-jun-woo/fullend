//ff:type feature=ssac-gen type=model
//ff:what 변수의 출처 정보(DDL/func)를 나타내는 구조체
package generator

// varSource는 변수의 출처 정보를 나타낸다.
type varSource struct {
	Kind      string // "ddl" or "func"
	ModelName string // DDL: "Workflow", Func: "CheckCredits"
}
