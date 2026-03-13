# Phase012: shorthand @response varName 교차 검증 ✅ 완료

## 목표

`@response varName` (shorthand) 사용 시 변수 타입의 JSON 필드명과 OpenAPI 응답 스키마 property를 비교하여 불일치를 검출한다.

## 배경

현재 `checkResponseFields`는 shorthand `@response varName`을 건너뛴다 (`extractResponseFieldKeys`가 `seq.Target != ""`이면 `nil` 반환). 명시적 필드(`@response { field: var }`)만 검증.

결과: `@call auth.IssueTokenResponse token = auth.IssueToken(...)` → `@response token`에서 `IssueTokenResponse.AccessToken`(json tag: `access_token`)과 OpenAPI `AccessToken`의 불일치를 못 잡음.

## 변수 타입 추적 경로

```
@response token
    ↓ seq.Target = "token"
    ↓ 시퀀스 역추적: token은 어디서 할당?
@call auth.IssueTokenResponse token = auth.IssueToken(...)
    ↓ seq.Result.Type = "IssueTokenResponse"
    ↓ funcspec에서 ResponseFields 조회
    ↓ [{Name: "AccessToken", Type: "string"}]
    ↓ json 태그 확인 → "access_token"
    ↓ OpenAPI 응답 property: "AccessToken"
    ↓ "access_token" ≠ "AccessToken" → ERROR
```

## 변수 출처별 타입 해석

| 출처 | 타입 결정 | 필드 소스 |
|---|---|---|
| `@call` | `seq.Result.Type` → funcspec `ResponseFields` | Go struct (json 태그) |
| `@get` / `@put` / `@post` / `@delete` | `seq.Result.Type` → DDL 테이블 컬럼 | symbol table columns |
| `@state` | 없음 (상태 문자열) | 해당 없음 |

### `@call` 결과 타입
funcspec에서 `ResponseFields`를 조회. **json 태그가 있으면 json 태그를, 없으면 Go 필드명을 JSON 키로 사용.**

### `@get`/`@put` 등 모델 결과 타입
DDL symbol table에서 컬럼 목록 조회. 컬럼명은 snake_case → 이미 JSON 키와 동일.

## 변경 내용

### 1. `internal/funcspec/parser.go` — json 태그 파싱 추가

`Field` struct에 `JSONName string` 필드 추가. Go AST에서 struct tag를 읽어 `json:"xxx"` 값을 추출.

```go
type Field struct {
    Name     string // Go struct field name
    Type     string // Go type
    JSONName string // json tag name (empty = use Name)
}
```

파서에서 `ast.Field.Tag`를 읽어 `json:"xxx"` 값 추출:
```go
if field.Tag != nil {
    tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
    if jn, ok := tag.Lookup("json"); ok {
        jn = strings.Split(jn, ",")[0] // omitempty 등 제거
        if jn != "" && jn != "-" {
            f.JSONName = jn
        }
    }
}
```

### 2. `internal/crosscheck/ssac_openapi.go` — shorthand 검증 추가

`checkResponseFields`에서 shorthand 처리 분기 추가:

```go
if responseFields == nil {
    // shorthand: @response varName — 변수 타입에서 필드 추출
    shorthandFields := resolveShorthandResponseFields(fn, funcSpecs, st)
    if shorthandFields == nil {
        continue // 타입 추적 불가 → 스킵
    }
    // OpenAPI property와 비교
    for _, jf := range shorthandFields {
        if !opProps[jf] {
            // ERROR
        }
    }
    for prop := range opProps {
        if !shorthandSet[prop] {
            // WARNING
        }
    }
    continue
}
```

### 3. `internal/crosscheck/ssac_openapi.go` — `resolveShorthandResponseFields` 신규

```go
// resolveShorthandResponseFields는 @response varName의 변수 타입을 추적하여
// JSON 필드명 목록을 반환한다.
func resolveShorthandResponseFields(
    fn ssacparser.ServiceFunc,
    funcSpecs []funcspec.FuncSpec,
    st *ssacvalidator.SymbolTable,
) []string
```

로직:
1. `@response`의 `seq.Target`으로 변수명 획득
2. 시퀀스 역추적: `seq.Result.Var == varName`인 시퀀스 찾기
3. **Wrapper 타입 스킵**: `seq.Result.Wrapper != ""`(Page, Cursor)이면 `nil` 반환 — wrapper 구조(`items`/`total`)는 고정이라 json 태그 불일치 위험 없음
4. `@call` → funcspec에서 `ResponseFields` 조회, `JSONName` 우선 사용
5. `@get`/`@put`/`@post`/`@delete` → DDL symbol table에서 컬럼명 조회

### 4. `internal/crosscheck/crosscheck.go` — `CrossValidateInput`에 funcspec 전달

`CheckSSaCOpenAPI` 호출 시 funcspec 정보가 필요. `CrossValidateInput`에 이미 `FullendPkgSpecs`, `ProjectFuncSpecs`가 있으므로, `CheckSSaCOpenAPI` 시그니처에 funcspec 슬라이스를 추가.

### 5. `internal/funcspec/parser_test.go` — json 태그 파싱 테스트

`IssueTokenResponse { AccessToken string \`json:"access_token"\` }` 같은 struct에서 `JSONName = "access_token"`이 추출되는지 검증.

### 6. `internal/crosscheck/ssac_openapi_test.go` — shorthand 검증 테스트

- `@response token` + funcspec `AccessToken(json:access_token)` + OpenAPI `AccessToken` → ERROR
- `@response token` + funcspec `AccessToken(json:access_token)` + OpenAPI `access_token` → 통과
- `@response user` + DDL columns `[id, email, name]` + OpenAPI `[id, email, name]` → 통과

## 영향 없는 범위

- 명시적 `@response { field: var }` — 기존 로직 유지
- `@response` 없는 함수 — 변경 없음
- SSaC/STML 코드 — 변경 없음

## 검증

```
go test ./internal/funcspec/...
go test ./internal/crosscheck/...
go test ./...
```

전체 통과 확인. gigbridge validate에서 Login `AccessToken` 불일치 검출 확인.
