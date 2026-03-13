# Phase017: 파싱 통합 — 1회 파싱, N회 공유 ✅ 완료

## 목표

validate, gen, status, chain 4개 경로에서 동일 SSOT를 반복 파싱하는 문제를 해결한다. 한 번 파싱한 결과를 `ParsedSSOTs` 구조체에 캐시하고, 모든 소비자에게 포인터로 전달한다.

## 현황: 중복 파싱 전수 조사

| SSOT | 현재 파싱 횟수 | 위치 |
|------|-------------|------|
| OpenAPI | 5 | validate.go, gen.go(genGlue), status.go, chain.go, genmodel |
| DDL (SymbolTable) | 5 | validate.go, gen.go(genSSaC), gen.go(genGlue), status.go, chain.go |
| fullend.yaml | 3 | validate.go, gen.go(genGlue), gen.go(determineModulePath) |
| SSaC | 5 | validate.go, gen.go(genSSaC), gen.go(genGlue), status.go, chain.go |
| STML | 3 | validate.go, gen.go(genSTML), status.go |
| Func (pkg+project) | 7 | validate.go×2, gen.go(injectFuncErrStatus)×2, gen.go(genFunc), chain.go×2 |
| States | 5 | validate.go, gen.go(genGlue), gen.go(genStateMachines), status.go, chain.go |
| Policy | 5 | validate.go, gen.go(genGlue), gen.go(genAuthz), status.go, chain.go |
| **합계** | **38** | — |

통합 후 각 SSOT는 호출(invocation)당 **1회**만 파싱. 38 → 8.

## 안전성 확인

| 컴포넌트 | 파싱 결과 변형(mutate)? | 공유 안전? |
|----------|----------------------|----------|
| ssac/validator (SymbolTable 구축) | 자체 st.Models, st.Operations 구축 | ✅ 원본 DDL/OpenAPI 안 건드림 |
| stml/validator (SymbolTable 구축) | 자체 st.Operations 구축 | ✅ 원본 OpenAPI 안 건드림 |
| crosscheck/* | 완전 읽기 전용 | ✅ |
| gluegen/* | 읽기 전용 | ✅ |
| injectFuncErrStatus | st.Models에 쓰기 | ⚠️ gen 전용 — validate의 st와 분리 필요 |

**주의**: `injectFuncErrStatus`는 SymbolTable을 mutate한다. gen 경로에서는 별도 SymbolTable을 사용하거나, validate 완료 후 gen용으로 재로드해야 한다. → gen에서만 `LoadSymbolTable` 1회 추가 호출 (validate의 st를 오염시키지 않기 위해).

## 설계

### ParsedSSOTs 구조체

```go
// internal/orchestrator/parsed.go
package orchestrator

type ParsedSSOTs struct {
    Config       *projectconfig.ProjectConfig // fullend.yaml
    OpenAPIDoc   *openapi3.T                  // api/openapi.yaml
    SymbolTable  *ssacvalidator.SymbolTable    // DDL + sqlc + OpenAPI 심볼
    ServiceFuncs []ssacparser.ServiceFunc      // service/**/*.ssac
    STMLPages    []stmlparser.Page             // frontend/*.html
    States       []*statemachine.StateDiagram  // states/*.md
    Policies     []*policy.Policy              // policy/*.rego
    FuncSpecs    []funcspec.FuncSpec           // func/**/*.go (project)
    PkgFuncSpecs []funcspec.FuncSpec           // pkg/**/*.go (fullend 내장)
    HurlFiles    []string                      // tests/*.hurl
    ModelDir     string                        // model/ 경로
}
```

### ParseAll 함수

```go
// internal/orchestrator/parsed.go
func ParseAll(root string, detected []DetectedSSOT, skip map[SSOTKind]bool) (*ParsedSSOTs, *reporter.Report)
```

- `detected`에 있는 SSOT를 순회하며 각 파서를 **1회** 호출
- 파싱 에러는 `reporter.Report`에 기록
- skip된 SSOT는 파싱하지 않음

### 소비자 변경

| 함수 | 변경 전 | 변경 후 |
|------|--------|--------|
| `Validate()` | 내부에서 각 SSOT 직접 파싱 | `ParseAll()` 호출 → `parsed.*` 사용 |
| `GenWith()` | Validate 후 다시 파싱 | Validate에서 받은 `parsed` 재사용 (단, SymbolTable은 gen용 별도) |
| `Status()` | 독립 파싱 | `ParseAll()` 호출 → `parsed.*` 사용 |
| `Chain()` | 독립 파싱 | `ParseAll()` 호출 → `parsed.*` 사용 |
| `genGlue()` | OpenAPI, SSaC, SymbolTable, States, Policy 재파싱 | `parsed.*`에서 전달받음 |
| `genSSaC()` | SSaC + SymbolTable 재파싱 | `parsed.*`에서 전달받음, gen용 st 별도 복사 |
| `genStateMachines()` | States 재파싱 | `parsed.*`에서 전달받음 |
| `genAuthz()` | Policy 재파싱 | `parsed.*`에서 전달받음 |

### gen 경로 SymbolTable 처리

```go
// gen에서는 validate의 st를 오염시키지 않기 위해 별도 로드
// injectFuncErrStatus가 st.Models를 mutate하므로
genST, _ := ssacvalidator.LoadSymbolTable(specsDir)  // gen 전용 1회
injectFuncErrStatus(genST, specsDir)
```

또는 SymbolTable에 `Clone()` 메서드를 추가하여 shallow copy 후 mutate:

```go
genST := parsed.SymbolTable.Clone()
injectFuncErrStatus(genST, specsDir)
```

→ `Clone()` 방식 채택 (파싱 0회 추가, 메모리 복사만).

## 변경 계획

### 1단계: ParsedSSOTs + ParseAll 신규

- `internal/orchestrator/parsed.go` 신규 생성
- `ParsedSSOTs` 구조체 + `ParseAll()` 함수

### 2단계: Validate 리팩터

- `validate.go`: `Validate()` 시그니처 변경
  - 내부 파싱 로직 제거 → `ParseAll()` 결과 사용
  - 개별 validate 함수(validateOpenAPI 등)는 파싱 결과를 인자로 받도록 변경
  - crosscheck 호출부는 이미 파싱 결과를 전달하므로 변경 최소

### 3단계: Gen 리팩터

- `gen.go`: `GenWith()` 변경
  - Validate 단계에서 받은 `parsed`를 gen 단계로 전달
  - `genSSaC()`, `genGlue()`, `genStateMachines()`, `genAuthz()` 시그니처 변경
  - SymbolTable은 `Clone()` 후 `injectFuncErrStatus` 적용

### 4단계: Status/Chain 리팩터

- `status.go`: `Status()` → `ParseAll()` 사용
- `chain.go`: `Chain()` → `ParseAll()` 사용

### 5단계: SymbolTable.Clone() 추가

- `internal/ssac/validator/symbol.go`에 `Clone()` 메서드 추가
- Models map shallow copy (injectFuncErrStatus가 Models만 mutate)

### 6단계: genmodel 정리

- `internal/genmodel/genmodel.go`의 OpenAPI 로드를 외부에서 주입받도록 변경 (선택)

## 변경 파일 요약

| 범위 | 파일 수 | 변경 |
|------|--------|------|
| parsed.go | 1 | 신규 |
| validate.go | 1 | ParseAll 사용으로 리팩터 |
| gen.go | 1 | parsed 전달, gen 함수 시그니처 변경 |
| status.go | 1 | ParseAll 사용으로 리팩터 |
| chain.go | 1 | ParseAll 사용으로 리팩터 |
| ssac/validator/symbol.go | 1 | Clone() 메서드 추가 |
| genmodel/genmodel.go | 1 | OpenAPI doc 주입 (선택) |
| **합계** | **~7** | — |

## 사전 검증: 파싱 멱등성 테스트 ✅ 완료

통합 전 전제 조건인 "같은 입력 → 같은 결과"를 gigbridge 실 데이터로 검증했다.

**테스트**: `internal/orchestrator/parse_idempotent_test.go`
**방법**: 각 파서를 동일 경로로 2회 호출 → `reflect.DeepEqual` 비교

| SSOT | 파서 | 결과 |
|------|------|------|
| OpenAPI | `openapi3.NewLoader().LoadFromFile()` | ✅ 동일 |
| DDL (SymbolTable) | `ssacvalidator.LoadSymbolTable()` | ✅ 동일 |
| SSaC | `ssacparser.ParseDir()` | ✅ 동일 |
| STML | `stmlparser.ParseDir()` | ✅ 동일 |
| States | `statemachine.ParseDir()` | ✅ 동일 |
| Policy | `policy.ParseDir()` | ✅ 동일 |
| FuncSpec | `funcspec.ParseDir()` | ✅ 동일 |
| ProjectConfig | `projectconfig.Load()` | ✅ 동일 |

```
=== RUN   TestParseIdempotency
--- PASS: TestParseIdempotency (0.01s)
    --- PASS: TestParseIdempotency/OpenAPI (0.01s)
    --- PASS: TestParseIdempotency/SymbolTable (0.00s)
    --- PASS: TestParseIdempotency/SSaC (0.00s)
    --- PASS: TestParseIdempotency/STML (0.00s)
    --- PASS: TestParseIdempotency/States (0.00s)
    --- PASS: TestParseIdempotency/Policy (0.00s)
    --- PASS: TestParseIdempotency/FuncSpec (0.00s)
    --- PASS: TestParseIdempotency/ProjectConfig (0.00s)
PASS
ok  	github.com/geul-org/fullend/internal/orchestrator	0.016s
```

**결론**: 8개 파서 모두 멱등성 확인. 1회 파싱 후 결과를 공유해도 안전하다.

## 구현 후 검증

```bash
go build ./cmd/fullend/
go test ./...
./fullend validate specs/gigbridge
./fullend gen specs/gigbridge artifacts/gigbridge
./fullend status specs/gigbridge
cd artifacts/gigbridge/backend && go build -o server ./cmd/
```

## 리스크

| 리스크 | 대응 |
|--------|------|
| validate 함수 시그니처 변경 범위 | 내부 함수만 변경, 외부 API(Validate/Gen/Status) 시그니처 유지 |
| SymbolTable Clone 누락 필드 | gen 테스트로 즉시 발견 |
| genmodel 외부 주입 누락 | 선택 사항, 후순위 가능 |
| injectFuncErrStatus 외 다른 mutator 존재 | 전수 조사 완료, 현재 유일 |
