# internal/gen

파싱된 SSOT로부터 백엔드 + 프론트엔드 + Hurl 스모크 테스트 산출물 전체를 생성하는 최상위 코드젠 오케스트레이터.

## 용어 주석

본 README에서 **feature**는 `specs/service/<폴더>/*.ssac`의 서브폴더 이름을 뜻한다.
예: `specs/service/auth/login.ssac` → feature = `auth`.
내부 코드 식별자는 역사적 이유로 `Domain`이라는 이름을 사용한다 (`ServiceFunc.Domain`, `hasDomains()`, `uniqueDomains()` 등). 읽을 때 feature로 치환해도 의미는 동일.

## 진입점

```go
// internal/gen/generate.go
func Generate(parsed *genapi.ParsedSSOTs, cfg *genapi.GenConfig, stmlOut *genapi.STMLGenOutput) error
```

호출 순서 (`generate.go:12-26`):

```
1. selectBackend(parsed.Config).Generate(parsed, cfg)
     → internal/gen/gogin (현재 유일한 백엔드 구현)
2. react.Generate(parsed, cfg, stmlOut)
     → React 글루 (App, main, api 클라이언트, 정적 파일)
3. hurl.Generate(parsed, cfg)
     → tests/smoke.hurl
```

이 세 단계는 **순차 실행**이며 실패 시 즉시 반환한다.
React와 Hurl은 OpenAPI contract 기반이라 백엔드 구현과 무관하게 돌아간다 (`generate.go:18, 22` 주석).

## 입력 타입

`*genapi.ParsedSSOTs` — 파싱된 모든 SSOT 집합.
| 필드 | 용도 |
|------|------|
| `OpenAPIDoc *openapi3.T` | 라우트/스키마/보안 스펙 |
| `ServiceFuncs []ssacparser.ServiceFunc` | SSaC 함수 정의 |
| `Policies []*policy.Policy` | OPA Rego 정책 |
| `StateDiagrams []statemachine.StateDiagram` | Mermaid 상태 머신 |
| `Config *projectconfig.ProjectConfig` | fullend.yaml |

`*genapi.GenConfig` — 생성 타깃 설정.
| 필드 | 용도 |
|------|------|
| `ArtifactsDir` | 산출물 루트 경로 |
| `SpecsDir` | DDL/Query 등 원천 스펙 경로 |
| `ModulePath` | 생성되는 Go module 경로 |

`*genapi.STMLGenOutput` — STML 생성기가 미리 만들어둔 결과.
| 필드 | 용도 |
|------|------|
| `Deps map[string]string` | 페이지가 요구하는 npm 패키지 |
| `Pages []string` | 생성된 페이지 파일 이름 |
| `PageOps map[string]string` | page → primary operationID 매핑 |

## 백엔드 선택

```go
// internal/gen/select_backend.go
func selectBackend(cfg *projectconfig.ProjectConfig) genapi.Backend {
    return &gogin.GoGin{}
}
```

현재 분기 없이 `gogin` 고정. 향후 `cfg.Backend` 필드로 분기할 것을 주석으로 남김.

## 서브 패키지

| 패키지 | 책임 |
|--------|------|
| `internal/gen/gogin/` | Go + Gin 백엔드 (`backend/cmd/main.go`, `backend/internal/{service,model,auth,middleware,authz}/`) |
| `internal/gen/react/` | React 프론트엔드 글루 (`frontend/src/{App,main,api}.tsx`, 정적 설정) |
| `internal/gen/hurl/` | Hurl 스모크 테스트 (`tests/smoke.hurl`) |

각 서브 패키지의 상세는 하위 README 참조.

## 실행 순서의 의미

- **Backend 먼저** — 모델/핸들러/Server 구조 확립. 이후 단계의 contract(OpenAPI)는 이미 파싱 완료 상태.
- **React 다음** — STML이 `frontend/src/pages/*.tsx`를 **선행 생성**한 상태여야 `scanPageFiles()`가 라우트를 조립한다. 상위 오케스트레이터(`internal/orchestrator/gen_stml.go`)가 STML 단계를 먼저 실행.
- **Hurl 마지막** — OpenAPI + 상태 머신 + Policy 기반이라 백엔드/프론트엔드 산출물에 의존하지 않음. 순서는 논리적 분류 목적.
