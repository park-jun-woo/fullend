# Phase049: Funcspec 파서 패키지 레벨 타입 해석 수정

## 목표

`internal/funcspec` 파서가 `@func` 파일만 AST 파싱하여 같은 패키지의 별도 파일에 있는 Request/Response 구조체를 인식하지 못하는 버그를 수정한다.

## 배경

filefunc 도입(Phase031~043)으로 pkg/auth, pkg/mail의 파일이 분리되었다:
- `hash_password.go` → `//ff:func` (@func + 함수 본체)
- `hash_password_request.go` → `//ff:type` (Request 구조체)
- `hash_password_response.go` → `//ff:type` (Response 구조체)

`ParseFile`은 `@func`이 있는 파일 하나만 `go/parser`로 읽으므로, 분리된 Request/Response 필드를 "0개"로 보고한다.

### 영향

- `fullend validate` 교차 검증(Cross) 단계에서 `Func ↔ SSaC` 체크 실패
- 영향받는 built-in 패키지: `pkg/auth` (hashPassword, verifyPassword), `pkg/mail` (sendTemplateEmail)
- `fullend gen` 진행 불가 (validate 실패 시 gen 차단)
- gigbridge-try03 더미 테스트에서 발견

## 근본 원인

`internal/funcspec/parse_file.go:43` — `for _, decl := range f.Decls` 루프가 단일 파일의 AST만 순회. `ParseFile`은 "단일 파일 파서"로서 역할이 맞으므로, 패키지 레벨 타입 보충은 `ParseDir`에서 처리해야 한다.

테스트(`parser_test.go`)는 Request/Response가 같은 파일에 있는 케이스만 검증하여 이 패턴을 놓쳤다.

## 변경 파일 목록

### 1. ParseDir에서 패키지 레벨 타입 보충

- **파일**: `internal/funcspec/parse_dir.go`
- **변경**: 디렉토리 내 모든 `.go` 파일을 먼저 한 번 파싱하여 타입 맵(구조체명 → []Field)을 구성한 뒤, `@func`이 있는 FuncSpec의 `RequestFields`/`ResponseFields`가 비어 있으면 타입 맵에서 보충
- **`parse_file.go`는 수정하지 않음** — 단일 파일 파서 책임 유지

```
ParseDir 수정 흐름:
1. 디렉토리별로 모든 .go 파일을 go/parser로 파싱 → 구조체명 → []Field 맵 구성
2. @func이 있는 FuncSpec에 RequestFields가 비어 있으면 맵에서 expectedRequest 조회하여 보충
3. ResponseFields도 동일
```

### 2. built-in 패키지 override 차단

- **파일**: `internal/orchestrator/validate_funcspec.go`
- **변경**: 기존 `validateFunc()` 함수에 built-in override 감지 로직 추가
- **방법**: `FullendPkgSpecs`에서 패키지명을 동적 추출하여 `ProjectFuncSpecs`와 비교. 겹치면 ERROR
- **하드코딩하지 않음** — `FullendPkgSpecs`가 SSOT

```go
// validate_funcspec.go 내 추가 로직
builtinPkgs := map[string]bool{}
for _, s := range fullendPkgSpecs {
    builtinPkgs[s.Package] = true
}
for _, s := range projectFuncSpecs {
    if builtinPkgs[s.Package] {
        // ERROR: func/<pkg>: built-in 패키지 "<pkg>"를 override할 수 없습니다
    }
}
```

- **함수 시그니처 변경**: `validateFunc(specs []funcspec.FuncSpec)` → `validateFunc(projectSpecs, fullendSpecs []funcspec.FuncSpec)`
- **호출부 수정**: `internal/orchestrator/validate_with.go:65` — `validateFunc(parsed.ProjectFuncSpecs, parsed.FullendPkgSpecs)` 로 변경
- **에러 메시지**: `func/auth: built-in 패키지 "auth"를 override할 수 없습니다. 커스텀 패키지명을 사용하세요 (예: func/myauth/)`
- **사유**:
  - codegen의 reexport.go와 반드시 redeclare 충돌 — 구조상 회피 불가
  - built-in의 검증된 구현(bcrypt 해싱 등)을 미검증 코드로 대체하는 보안 위험
  - 같은 `auth.HashPassword`가 프로젝트마다 다르게 동작하는 디버깅 난이도 상승

### 3. 테스트 추가

- **파일**: `internal/funcspec/parser_test.go`
- **추가 케이스**:
  - `TestParseDirSplitFiles` — `@func` 파일과 Request/Response 구조체가 별도 파일에 있는 케이스. ParseDir 결과의 RequestFields/ResponseFields가 정상 채워지는지 검증
- **파일**: `internal/orchestrator/validate_funcspec_test.go` (신규 또는 기존에 추가)
- **추가 케이스**:
  - built-in 패키지명과 겹치는 ProjectFuncSpecs가 있으면 ERROR 반환 확인

## 수정하지 않는 파일

| 파일 | 사유 |
|---|---|
| `internal/funcspec/parse_file.go` | 단일 파일 파서 책임 유지. 패키지 레벨 보충은 ParseDir의 역할 |
| `internal/orchestrator/gen_func.go` (func-gen) | validate에서 override를 차단하므로 gen까지 도달하지 않음. 별도 방어 불필요 |
| `internal/crosscheck/build_func_spec_map.go` | project override 자체가 차단되므로 기존 "project overrides fullend" 로직은 dead code가 되지만, 삭제는 별도 Phase에서 판단 |
| `internal/orchestrator/scan_func_imports.go` | fullend built-in import skip 로직 정상 동작 |

## 의존성

없음

## 검증 방법

1. `go test ./internal/funcspec/...` — `TestParseDirSplitFiles` 통과
2. `fullend validate specs/dummys/zenflow-try05` — pkg/auth, pkg/mail 관련 `Func ↔ SSaC` 에러 0건
3. `fullend validate specs/dummys/gigbridge-try03` — 로컬 func/auth override 없이 통과
4. gigbridge-try03에 `func/auth/` 추가 후 `fullend validate` → built-in override ERROR 발생 확인
5. `go test ./...` — 전체 테스트 통과
