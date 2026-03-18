# Phase002 — internal/policy 파서 교체: Rego regex → Go AST + YAML

## 목표

`internal/policy/`의 Rego regex 파서를 Go AST + toulmin YAML 그래프 파서로 교체한다. 출력 데이터 구조(`Policy`, `AllowRule`, `OwnershipMapping`)는 유지하여 downstream(crosscheck, codegen) 호환성을 보장한다.

## 배경

현재 `internal/policy/`는 15개 파일로 Rego 텍스트를 regex로 파싱한다. 전환 후 두 종류의 파일을 파싱한다:
1. `*.go` — Go AST로 함수 선언 + `//tm:backing` 어노테이션 추출
2. `*.yaml` — toulmin YAML 그래프 정의에서 role, qualifier, defeats 추출

## 현재 구조 (삭제/교체 대상)

```
internal/policy/
├── parse_dir.go                  # *.rego 파일 탐색 → *.go + *.yaml 탐색으로 변경
├── parse_file.go                 # Rego regex 파싱 → Go AST + YAML 파싱으로 교체
├── extract_allow_rules.go        # "allow if {" 분할 → 삭제
├── process_allow_block.go        # Rego 블록 내 패턴 매칭 → 삭제
├── parse_action_set.go           # 액션 셋 파싱 → 삭제
├── find_closing_brace.go         # 중괄호 매칭 → 삭제
├── parse_ownership_line.go       # @ownership 주석 파싱 → Go 구조체 파싱으로 변경
├── policy_type.go                # Policy 구조체 → 유지
├── allow_rule.go                 # AllowRule 구조체 → 유지 + 필드 추가
├── ownership_mapping.go          # OwnershipMapping 구조체 → 유지
├── policy_action_resource_pairs.go  # 유지
├── policy_resources_using_owner.go  # 유지
├── policy_ownership_for.go          # 유지
```

## 새 파서 설계

### ParseDir 변경

```go
// 기존: *.rego 파일 탐색
// 신규: authz/ 디렉토리에서 *.go + *.yaml 파싱
func ParseDir(dir string) ([]*Policy, error) {
    // 1. *.yaml 파싱 → GraphDef (rules, defeats)
    // 2. *.go 파싱 → 함수 선언 + //tm:backing + 본문 AST 분석
    // 3. GraphDef + Go 함수 정보를 결합하여 Policy 구조체 생성
}
```

### Go 파일 파싱

```go
func parseGoFile(path string) ([]parsedFunc, []OwnershipMapping, error) {
    fset := token.NewFileSet()
    f, _ := parser.ParseFile(fset, path, nil, parser.ParseComments)

    for _, decl := range f.Decls {
        fn, ok := decl.(*ast.FuncDecl)
        // //tm:backing 주석에서 backing 추출
        // 함수 본문 AST에서 claims/role/owner 참조 탐지
    }
    // var Ownerships 선언에서 OwnershipMapping 추출
}
```

### YAML 파일 파싱

```go
func parseYAMLGraph(path string) (*graphDef, error) {
    // toulmin의 internal/graphdef 패키지 재사용 또는 동일 구조 파싱
    // rules: name, role, qualifier
    // defeats: from, to
}
```

### Go AST + YAML 결합

```go
type parsedFunc struct {
    Name       string
    Backing    string
    UsesClaims bool     // 본문에서 .Claims 필드 접근 여부
    UsesOwner  bool     // 본문에서 .ResourceOwnerID 접근 여부
    UsesRole   bool     // 본문에서 .Claims.Role 접근 여부
    RoleValue  string   // Role 비교값 (있는 경우)
    ClaimsRefs []string // 참조하는 claims 필드 목록
}

// YAML의 role/defeats 정보 + Go의 본문 분석 결과 → AllowRule
```

### AllowRule 확장

```go
type AllowRule struct {
    Actions    []string
    Resource   string
    UsesOwner  bool
    UsesRole   bool
    UsesClaims bool     // 신규: claims 참조 여부 (TODO007 D1 해소)
    RoleValue  string
    SourceLine int
    // toulmin 메타데이터
    Role       string   // warrant, rebuttal, defeater (YAML에서)
    Qualifier  float64  // YAML에서 (기본 1.0)
    Defeats    []string // YAML에서
    Backing    string   // //tm:backing에서
}
```

### Claims/Role/Owner 참조 탐지

Go AST에서 함수 본문을 순회하여:
- `c.Claims.UserID` → UsesClaims = true, ClaimsRefs에 "user_id" 추가
- `c.Claims.Role` → UsesRole = true
- `c.ResourceOwnerID` → UsesOwner = true

### Ownership 파싱

```go
// 기존: Rego 주석 `# @ownership workflow: workflows.org_id` regex 파싱
// 신규: Go 변수 선언 `var Ownerships = []OwnershipMapping{...}` AST 파싱
func parseOwnerships(f *ast.File) []OwnershipMapping
```

## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| `internal/policy/parse_dir.go` | *.rego → *.go + *.yaml 탐색 |
| `internal/policy/parse_file.go` | 전면 교체: Go AST + YAML 파싱 |
| `internal/policy/extract_allow_rules.go` | 삭제 |
| `internal/policy/process_allow_block.go` | 삭제 |
| `internal/policy/parse_action_set.go` | 삭제 |
| `internal/policy/find_closing_brace.go` | 삭제 |
| `internal/policy/parse_ownership_line.go` | Go 구조체 파싱으로 교체 |
| `internal/policy/parse_yaml_graph.go` (신규) | YAML 그래프 파싱 |
| `internal/policy/parse_go_funcs.go` (신규) | Go AST 함수 파싱 + //tm:backing 추출 |
| `internal/policy/merge_graph_funcs.go` (신규) | YAML + Go 함수 결합 → AllowRule 생성 |
| `internal/policy/allow_rule.go` | Role, Qualifier, Defeats, Backing, UsesClaims 필드 추가 |
| `internal/policy/parser_test.go` | 테스트 전면 교체 |
| `internal/policy/parser_claims_test.go` | 테스트 교체 |

## 의존성

- `go/ast`, `go/parser`, `go/token` — Go 표준 라이브러리
- `gopkg.in/yaml.v3` — YAML 파싱 (이미 fullend 의존성에 포함 가능)
- `github.com/park-jun-woo/toulmin/pkg/toulmin` — ParseAnnotation 함수 (//tm: 파싱)

## 검증 방법

- 기존 `parser_test.go` 테스트를 새 포맷으로 변환하여 동일 출력 확인
- zenflow-try06의 새 authz/ 파일을 파싱하여 기존 Rego 파싱 결과와 동일한 Policy 구조체 생성 확인
- `go test ./internal/policy/...` 통과
