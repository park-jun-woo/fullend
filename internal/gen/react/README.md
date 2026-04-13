# internal/gen/react

React 프론트엔드의 **글루 코드**를 생성한다.
페이지 컴포넌트(`src/pages/*.tsx`) 자체는 이 패키지가 아니라 `internal/stml/generator/`가 만든다.
이 패키지는 그 페이지들을 감싸는 **App, main, api client, 빌드 설정** 을 담당.

## 진입점

```go
// internal/gen/react/generate.go
func Generate(parsed *genapi.ParsedSSOTs, cfg *genapi.GenConfig, stmlOut *genapi.STMLGenOutput) error
```

`generateFrontendSetup()` `generate_frontend_setup.go:14-42` 가 파일 쓰기를 일괄 수행.

## 아키텍처

```
STML PageSpec          ──→  internal/stml/generator     ──→  frontend/src/pages/*.tsx
                            (페이지 컴포넌트)

OpenAPI + stmlOut      ──→  internal/gen/react           ──→  frontend/{
                                                               package.json,
                                                               vite.config.ts,
                                                               tsconfig.json,
                                                               index.html,
                                                               src/main.tsx,
                                                               src/App.tsx,
                                                               src/api.ts,
                                                             }
```

**중요 의존성**: `src/App.tsx` 생성 시 `scanPageFiles()` `scan_page_files.go:13` 로 **디스크의 실제 페이지 파일을 스캔**한다.
따라서 **STML 생성이 선행되어야** 라우트가 조립됨.
순서는 상위 오케스트레이터(`internal/orchestrator/gen_stml.go`)가 보장.

## 산출물 맵

### 완전 정적 파일

각각 소스 템플릿이 Go 문자열 리터럴로 박혀 있음.

| 파일 | 내용 | 생성자 |
|------|------|-------|
| `frontend/vite.config.ts` | React 플러그인, `/api` 프록시 `localhost:8080` | `write_vite_config.go:11` |
| `frontend/tsconfig.json` | ES2020, JSX react-jsx, strict | `write_ts_config.go:11` |
| `frontend/index.html` | root div + main.tsx 모듈 스크립트 | `write_index_html.go:11` |

### 조건부 정적

| 파일 | 분기 | 생성자 |
|------|------|-------|
| `frontend/package.json` | 기본 deps + `stmlOut.Deps` 머지 | `write_package_json.go:12` (lines 18-21 머지) |
| `frontend/src/main.tsx` | `stmlOut.Deps["@tanstack/react-query"]` 유무 | `write_main_tsx.go:12` |

`main.tsx` 분기:
- 있음 → `main_tsx_with_tanstack.go:6` — `QueryClientProvider` 포함.
- 없음 → `main_tsx_without_tanstack.go:6` — 순수 `BrowserRouter`.

### 동적 생성

| 파일 | 입력 | 생성자 |
|------|------|-------|
| `frontend/src/api.ts` | `OpenAPIDoc.Paths` | `write_api_client.go:15` |
| `frontend/src/App.tsx` | 디스크 페이지 파일 + OpenAPI + `stmlOut.PageOps` | `write_app_tsx.go:17` |

## API 클라이언트 (`src/api.ts`) 생성

### 수집

`collectEndpoints()` `collect_endpoints.go:13`:
- `doc.Paths.Map()` 순회.
- `appendOperations()` `append_operations.go:9` 로 method × operation 평탄화.
- endpoint 구조체: `{method, path, opID, pathParams}`.
- `operationID` 기준 알파벳 정렬.

### 각 엔드포인트 → async 함수

`writeEndpointFunc()` `write_endpoint_func.go:12` 가 다음을 출력:

```ts
async function listCourses(params?: Record<string, any>) {
  // GET: URLSearchParams로 qs 조립, path params는 템플릿 리터럴로
  // POST/PUT/DELETE: path params 제외 후 JSON body
  const res = await fetch(`${BASE}${path}`, { method, headers, body })
  return res.json()
}
```

- 함수명: `lcFirst(operationID)` — `"ListCourses"` → `"listCourses"`.
- 경로 변환: `openAPIPathToTemplateLiteral()` `openapi_path_to_template_literal.go:9` — `/courses/{CourseID}` → `` `/courses/${courseID}` ``.
- Path param 추출: `extractPathParams()` `extract_path_params.go:9`.

### GET vs Mutation 분기

| 메서드 | 함수 | 동작 |
|-------|------|------|
| GET | `writeGetBody()` `write_get_body.go:12` | `URLSearchParams` 로 qs, path params 제외한 나머지 |
| POST/PUT/DELETE | `writeMutationBody()` `write_mutation_body.go:12` | path params 제외한 나머지를 JSON body, `Content-Type: application/json` |

### API namespace

`writeApiNamespace()` `write_api_namespace.go:12` 로 모든 함수를 export:

```ts
export const api = {
  ListCourses: listCourses,
  GetCourse:   getCourse,
  CreateCourse: createCourse,
  ...
}
```

### 엣지 케이스

- `doc == nil` → `export const api = {}` 공백 객체 (`write_api_client.go:19`).

## 라우트 생성 (`src/App.tsx`)

### 페이지 스캔

`scanPageFiles()` `scan_page_files.go:13`:
- `frontend/src/pages/*.tsx` 파일명 수집.
- 파일명(kebab-case) → 컴포넌트명(PascalCase) 변환.
- 디렉토리 부재 시 `nil` 반환 — 에러 아님 (`scan_page_files.go:15-17`).

### OpenAPI 매핑

`buildOpPaths()` `build_op_paths.go:9` — operationID → path 매핑.

### 라우트 조립

`buildAppRoutes()` `build_app_routes.go:8` 가 페이지마다 `resolveRoutePath()` 호출:

```
페이지 파일 → stmlPageOps[pageName] → opPaths[operationID] → route path
```

`resolveRoutePath()` `resolve_route_path.go:9`:
- 성공 경로: STML page → operationID → OpenAPI path → `openAPIPathToReactRoute()` 변환.
- Fallback: `/{page-name-kebab-case}`.

`openAPIPathToReactRoute()` `openapi_path_to_react_route.go:9`:
- `{PascalCase}` → `:camelCase`.
- 예: `/courses/{CourseID}` → `/courses/:courseID`.

`deduplicateRoutes()` `deduplicate_routes.go:9` — path 기준 중복 제거, 알파벳 정렬로 안정화.

### 결과

```tsx
<Routes>
  <Route path="/courses" element={<CourseListPage />} />
  <Route path="/courses/:courseID" element={<CourseDetailPage />} />
  ...
</Routes>
```

## 페이지 컴포넌트는 어디서?

이 패키지는 **페이지 컴포넌트 자체를 생성하지 않음**.
STML → TSX 변환은 `internal/stml/generator/react_target_generate_page.go:12` (`ReactTarget.GeneratePage`).

요약:
- `data-fetch` → `useQuery` 훅.
- `data-action` → `useMutation` + `useForm`.
- `data-bind="foo.bar"` → `{foo?.bar}` (null-safe chaining).
- `data-each="foo.items"` → `{items?.map(item => ...)}`.
- `data-state` → `useState` + 조건부 렌더.
- 중첩 fetch/action은 재귀 처리.

`@tanstack/react-query` 사용 여부는 STML 생성기가 `stmlOut.Deps` 에 기록하고, React 글루가 `main.tsx` 분기에 사용.

## 입력 의존성

| 소스 | 필드 | 용도 |
|------|------|------|
| `ParsedSSOTs` | `OpenAPIDoc *openapi3.T` | api.ts 엔드포인트, App.tsx path |
| `GenConfig` | `ArtifactsDir` | frontend/ 루트 |
| `STMLGenOutput` | `Deps map[string]string` | package.json 머지 + main.tsx 분기 |
| `STMLGenOutput` | `PageOps map[string]string` | App.tsx 라우트 매핑 |

`STMLGenOutput.Pages`는 현재 소비처 없음.

## 타입 생성 — 없음

- `components.schemas` → TypeScript interface 변환 **미실시**.
- API 함수 시그니처: `params?: Record<string, any>`, 반환 `any`.
- 타입 안전성은 호출자 책임.

향후 보완이 필요하면 `openapi-typescript` 같은 외부 도구 통합 또는 자체 schema → .d.ts 생성기 추가가 후보.

## 빠른 인덱스

| 기능 | 파일 | 함수 |
|------|------|------|
| 진입 | `generate.go` | `Generate:11` |
| 오케스트레이션 | `generate_frontend_setup.go` | `generateFrontendSetup:14` |
| package.json | `write_package_json.go` | `writePackageJSON:12` |
| main.tsx | `write_main_tsx.go` | `writeMainTSX:11` |
| api.ts | `write_api_client.go` | `writeAPIClient:15` |
| 엔드포인트 수집 | `collect_endpoints.go` | `collectEndpoints:13` |
| 엔드포인트 함수 | `write_endpoint_func.go` | `writeEndpointFunc:12` |
| api namespace | `write_api_namespace.go` | `writeApiNamespace:12` |
| App.tsx | `write_app_tsx.go` | `writeAppTSX:17` |
| 페이지 스캔 | `scan_page_files.go` | `scanPageFiles:13` |
| 라우트 조립 | `build_app_routes.go` | `buildAppRoutes:8` |
| 라우트 해석 | `resolve_route_path.go` | `resolveRoutePath:9` |
| 경로 변환 (API) | `openapi_path_to_template_literal.go` | `openAPIPathToTemplateLiteral:9` |
| 경로 변환 (Router) | `openapi_path_to_react_route.go` | `openAPIPathToReactRoute:9` |

## 특징 요약

- **이중 생성 구조** — 페이지는 STML generator, 글루는 이 패키지. 디스크 스캔으로 두 산출물이 연결됨.
- **contract-first** — OpenAPI가 유일한 API 형상 원천, 백엔드 구현과 무관하게 실행 가능.
- **정적 비중이 큼** — 7개 파일 중 4개가 (조건부 포함) 정적 템플릿. toulmin 같은 결정 그래프 필요 없음.
- **타입 안전성 공백** — TypeScript 타입 자동 생성은 없음.
- **기본 데이터 패칭 스택** — `fetch` + (선택) TanStack Query. Axios 등은 선택 의존성으로만.
