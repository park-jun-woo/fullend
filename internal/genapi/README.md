# internal/genapi

**SSOT 파싱 결과 통합 컨테이너**와 **백엔드 생성기 인터페이스**를 정의하는 타입 패키지.
로직은 없고, 다른 패키지가 주고받는 공유 타입만 둔다.

## 역할

- `orchestrator.ParseAll()` 이 **채우는** 자료형.
- `crosscheck`, `internal/gen/*` 가 **소비하는** 자료형.
- `Backend` 인터페이스로 백엔드 생성기(현재 Go+Gin)를 추상화.

전체 파일 4개, 약 60줄. 로직 없음.

## 파일 구성

| 파일 | 정의 | 역할 |
|------|------|------|
| `parsed_ssots.go` | `ParsedSSOTs` 구조체 | 모든 SSOT 파싱 결과 담는 컨테이너 |
| `gen_config.go` | `GenConfig` 구조체 | 출력 위치 / 모듈 경로 |
| `stml_gen_output.go` | `STMLGenOutput` 구조체 | STML 선행 생성 결과 |
| `backend.go` | `Backend` 인터페이스 | 백엔드 구현체 추상화 |

## ParsedSSOTs

`parsed_ssots.go:19-32`

```go
type ParsedSSOTs struct {
    Config           *projectconfig.ProjectConfig
    OpenAPIDoc       *openapi3.T
    SymbolTable      *ssacvalidator.SymbolTable
    ServiceFuncs     []ssacparser.ServiceFunc
    STMLPages        []stmlparser.PageSpec
    StateDiagrams    []*statemachine.StateDiagram
    Policies         []*policy.Policy
    ProjectFuncSpecs []funcspec.FuncSpec
    FullendPkgSpecs  []funcspec.FuncSpec
    HurlFiles        []string
    ModelDir         string
    StatesErr        error
}
```

| 필드 | 출처 파서 | 소비자 |
|------|----------|--------|
| `Config` | `projectconfig.Load` | gen_glue, gogin (auth/queue/session/...) |
| `OpenAPIDoc` | `openapi3.Loader.LoadFromFile` | gogin(라우팅/security), react(api client), hurl(scenario), crosscheck |
| `SymbolTable` | `ssacvalidator.LoadSymbolTable` | ssac generator, crosscheck |
| `ServiceFuncs` | `ssacparser.ParseDir` | gogin(handler), ssac generator |
| `STMLPages` | `stmlparser.ParseDir` | stml generator |
| `StateDiagrams` | `statemachine.ParseDir` | state-gen, hurl(BFS 순서) |
| `Policies` | `policy.ParseDir` | authz-gen, gogin(ownership), hurl(role 매핑) |
| `ProjectFuncSpecs` | `funcspec.ParseDir` | ssac generator(@error), validate/funcspec |
| `FullendPkgSpecs` | `funcspec.ParseDir(pkgRoot)` | validate/funcspec |
| `HurlFiles` | `filepath.Glob` | scenario 검증 |
| `ModelDir` | (경로만) | gogin |
| `StatesErr` | states 파싱 실패 시 보존 | validate/states |

States 파싱만 실패 내용을 필드에 보존하고, 나머지는 실패 시 nil.

## GenConfig

`gen_config.go:6-10`

```go
type GenConfig struct {
    ArtifactsDir string  // 산출물 루트 (backend/, frontend/, tests/ 생성됨)
    SpecsDir     string  // DDL/Query 등 원천 스펙 디렉토리
    ModulePath   string  // 생성되는 Go module 경로 (예: github.com/user/project)
}
```

- `ArtifactsDir` — `internal/gen/*` 전체가 이를 루트로 파일 출력.
- `SpecsDir` — DDL/query 파일을 직접 파싱하는 gogin에서 사용.
- `ModulePath` — 생성된 Go 파일의 import 경로에 삽입 (`gogin/generate.go:15-16`).

## STMLGenOutput

`stml_gen_output.go:7-11`

```go
type STMLGenOutput struct {
    Deps    map[string]string   // npm 패키지명 → 버전 범위 (예: "@tanstack/react-query" → "^5")
    Pages   []string            // 생성된 페이지 파일명
    PageOps map[string]string   // page명 → 대표 operationID
}
```

- `Deps` — react generator가 `package.json` 머지 + `main.tsx` 분기에 사용.
- `Pages` — 현재 소비처 없음. 향후 참조용.
- `PageOps` — react generator가 App.tsx 라우트 경로 결정에 사용.

## Backend 인터페이스

`backend.go:6-8`

```go
type Backend interface {
    Generate(parsed *ParsedSSOTs, cfg *GenConfig) error
}
```

구현체: `internal/gen/gogin.GoGin` (`gogin/generate.go:15`).
선택기: `internal/gen/select_backend.go:11` — 현재 분기 없이 `&gogin.GoGin{}` 고정 반환. 향후 `Config.Backend` 필드 기반 분기 예정 (주석으로 남겨둠).

## 데이터 흐름

```
orchestrator.GenWith()
  → ParseAll()                       genapi.ParsedSSOTs ← 채움
  → ValidateWith()                   ParsedSSOTs  ← 소비(검증)
  → runCodegenSteps()
      → genSTML()                    genapi.STMLGenOutput ← 반환
      → genGlue()
          → GenConfig 조립            genapi.GenConfig ← 채움
          → gen.Generate(parsed, cfg, stmlOut)
              → selectBackend().Generate(parsed, cfg)
              → react.Generate(parsed, cfg, stmlOut)
              → hurl.Generate(parsed, cfg)
```

## 외부 의존성

- `github.com/getkin/kin-openapi/openapi3` — `OpenAPIDoc *openapi3.T` 타입 제공.

내부 의존성(패키지 참조만):
`projectconfig`, `funcspec`, `policy`, `statemachine`, `ssacparser`, `ssacvalidator`, `stmlparser`.

## 설계 메모

- **순수 타입 패키지** — 행위 없음. 순환 의존 차단용.
- **컨테이너 방식** — nullable 필드 허용. `--skip`으로 빠진 SSOT는 nil로 남음.
- **확장 포인트** — `Backend` 인터페이스. 다른 프레임워크(gorilla/fiber/chi 등) 추가 시 구현체만 붙이고 `selectBackend` 분기.
