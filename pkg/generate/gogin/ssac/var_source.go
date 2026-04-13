//ff:type feature=ssac-gen type=model topic=type-resolve
//ff:what 변수의 출처 정보(DDL/func)를 나타내는 구조체
package ssac

// varSource는 변수의 출처 정보를 나타낸다.
type varSource struct {
	Kind      string // "ddl" or "func"
	ModelName string // DDL: "Workflow", Func: "CheckCredits"
}
