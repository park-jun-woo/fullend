# internal/stml/generator

STML(Structured Template Markup Language) 파싱 결과를 **React TSX 페이지 컴포넌트**로 변환.
`useQuery` / `useMutation` / `useForm` 훅과 JSX 렌더 트리를 모두 생성한다.

## 역할 경계

| 담당 | 이 패키지 | `internal/gen/react` |
|------|----------|---------------------|
| STML 페이지 → `src/pages/*.tsx` | **O** | — |
| `App.tsx` (라우팅) | — | **O** |
| `main.tsx` (provider) | — | **O** |
| `src/api.ts` (API 클라이언트) | — | **O** |
| `package.json`, vite/ts/html 설정 | — | **O** |
| 페이지 의존성(`@tanstack/react-query` 등) 집계 | **O** (`GenerateResult.Dependencies`) | 소비 |

이 패키지가 페이지 .tsx를 먼저 생성하면, `react` 패키지가 디스크에서 스캔해 라우팅/글루를 조립한다.

## 진입점

| 함수 | 파일 | 설명 |
|------|------|------|
| `Generate(pages, specsDir, outDir, opts…)` | `generate.go:8` | 기본 Target 사용 |
| `GenerateWith(target, pages, specsDir, outDir, opts…)` | `generate_with.go:14` | Target 명시 + 파일 쓰기 |
| `GeneratePage(page, specsDir, opts…)` | `generate_page.go:8` | 단일 페이지 TSX 문자열 |

### 시그니처

```go
func Generate(pages []parser.PageSpec, specsDir, outDir string, opts ...GenerateOptions) (*GenerateResult, error)
```

- `pages` — STML 파서가 만든 페이지 스펙.
- `specsDir` — 원본 STML 파일 경로.
- `outDir` — 산출 `.tsx` 출력 경로 (orchestrator가 `frontend/src/pages` 전달).
- `opts` — `APIImportPath`, `UseClient` 등.

### GenerateResult

`generate_result.go:6-9`

```go
type GenerateResult struct {
    Pages        int                 // 생성된 페이지 수
    Dependencies map[string]string   // npm 패키지 → 버전 범위
}
```

orchestrator `gen_stml.go:45-52` 가 이 결과 + pageName 목록 + pageOps(page → 대표 operationID) 를 `genapi.STMLGenOutput`으로 반환.

## 호출 경로

```
orchestrator.genSTML (orchestrator/gen_stml.go:16-58)
  └─ stmlgenerator.GenerateWith(profile.Frontend, pages, specsDir, outDir, opt)
      ├─ os.MkdirAll(outDir)
      ├─ for page in pages:
      │    ├─ code := target.GeneratePage(page, specsDir, opt)
      │    └─ os.WriteFile(outDir/{page.Name}.tsx, code)
      └─ return GenerateResult{Pages, Dependencies: target.Dependencies(pages)}
```

orchestrator는 `APIImportPath: "../api"`, `UseClient: false` 를 지정.

## Target 추상

`target_type.go:9-13`

```go
type Target interface {
    GeneratePage(page parser.PageSpec, specsDir string, opts GenerateOptions) string
    FileExtension() string
    Dependencies(pages []parser.PageSpec) map[string]string
}
```

구현: `ReactTarget` (`react_target_type.go:6-7`) — 현재 유일.

| 메서드 | 파일 |
|--------|------|
| `GeneratePage` | `react_target_generate_page.go:12-54` |
| `FileExtension` | `react_target_file_extension.go:5-7` → `".tsx"` |
| `Dependencies` | `react_target_dependencies.go:7-22` |
| `DefaultTarget` | `default_target.go:6-8` → ReactTarget 싱글톤 |

## 페이지 생성 파이프라인

`ReactTarget.GeneratePage` (`react_target_generate_page.go:12-54`):

```
1. collectImports(page, specsDir)               → importSet
2. collectFetchOps(page.Fetches)                 → fetch operationID[] (invalidation용)
3. collectAllActions(page.Children) + dedup      → 중첩 Action 전부 수집
4. renderImports(is, opt)                        → 'use client'?, import 블록
5. export default function <PascalName>() {
6.   renderPageHooks(page, is)                   → useParams, useQueryClient, useQuery들
7.   renderPageMutations(actions, fetchOps)      → useForm + useMutation들 (onSuccess에 invalidation)
8.   renderPageJSX(page)                         → return ( … )
9. }
```

## STML DSL → React 매핑

### data-fetch → useQuery

| 파일 | 함수 |
|------|------|
| `collect_imports.go:13-46` | Fetches 존재 시 `useQuery=true`, 파라미터에 `route.` 있으면 `useParams=true`, `Paginate/Sort/Filters` 존재 시 `useState=true` |
| `render_use_query.go:12-45` | `useQuery({ queryKey, queryFn: () => api.Op(args) })` 문자열 |
| `render_fetch_hooks.go:34` | useState(페이징/정렬/필터) + useQuery |
| `collect_fetch_ops.go:13` | 모든 Fetch OperationID 재귀 수집 (invalidation용) |

queryKey 조립:
- `'OperationID'` + route 파라미터들 + (있으면) `page, limit`, `sortBy, sortDir`, `filters`.

### data-action → useMutation

| 파일 | 함수 |
|------|------|
| `collect_action_imports.go:11-23` | Actions 존재 시 `useMutation=true`, `useQueryClient=true`; Fields 존재 시 `useForm=true` |
| `render_use_mutation.go:13-40` | `useMutation({ mutationFn, onSuccess: invalidate })` |
| `render_form_hook.go:14` | `const xForm = useForm()` |
| `render_page_mutations.go:12-19` | 모든 action의 form+mutation 렌더링 |

`mutationFn`:
- 폼 데이터 `data` + route 파라미터 병합해 `api.Op({...data, id})` 호출.

`onSuccess`:
- `fetchOps` 각각에 대해 `queryClient.invalidateQueries({ queryKey: [op] })` 를 생성.

### data-bind → JSX

`render_bind_jsx.go:12-17`:

```
{dataVar}.{fieldName} → <span>{dataVar.fieldName}</span>
```

- `dataVar` 는 컨텍스트별로 다름:
  - Fetch 자식 내부 → `"xxxData"` (예: `listReservationsData`)
  - Each 내부 → `"item"`
  - Action 폼 내부 → 폼 데이터/정적

### data-each → map

`render_each_jsx.go:13-38`:

```jsx
{dataVar.items?.map((item: any, index: number) => (
  <div key={index}>
    {/* Children: bind, static, each(중첩) … */}
  </div>
))}
```

null-safe chaining(`?.`) + 재귀 처리.

### data-state → 조건부 렌더

`render_state_jsx.go` — 상태값 기반 분기 렌더 (구체는 관련 파일 확인).

### 중첩 ChildNode 처리

`render_child_nodes.go:8-29` 가 dispatcher:

```
Kind → 렌더 함수
──────────────────
bind      → renderBindJSX
each      → renderEachJSX
state     → renderStateJSX
component → renderComponentJSX
static    → renderStaticJSX
action    → renderActionJSX
fetch     → renderFetchJSX (재귀)
```

## 훅 렌더링 상세

### renderPageHooks (`render_page_hooks.go:27`)

```
1. renderUseParams(is)                — const { id, name } = useParams()
2. const queryClient = useQueryClient() (if mutation)
3. renderFetchHooks(page.Fetches)     — useState + useQuery 각각
```

### renderPageMutations (`render_page_mutations.go:12-19`)

```
for each action in allActions:
    const xForm = useForm()          (if fields)
    const xMut  = useMutation({ mutationFn, onSuccess })
```

## JSX 구조 렌더링

| 함수 | 파일 | 역할 |
|------|------|------|
| `renderPageJSX` | `render_page_jsx.go:21` | `return ( ... )` 블록 |
| `renderPageJSXWithChildren` | `render_page_jsx_with_children.go:28` | 루트 정적 엘리먼트 + children |
| `renderPageJSXFallback` | `render_page_jsx_fallback.go:22` | children 없을 때 Fetch/Action 직접 나열 |
| `findRootElement` | `find_root_element.go` | children에서 루트 정적 요소 탐색 |
| `renderFetchJSX` | `render_fetch_jsx.go:31` | `{isLoading && …} {error && …} {data && …}` 3-way |
| `renderFetchJSXBody` | `render_fetch_jsx_body.go:33` | 필터·정렬·자식·페이지네이션 통합 |
| `renderFetchJSXFlatChildren` | `render_fetch_jsx_flat_children.go` | Fetch children 없을 때 기본 |
| `renderActionJSX` | `render_action_jsx.go:13` | Fields 유무 기준 form vs button |
| `renderActionForm` | `render_action_form.go:34` | `<form onSubmit={mut.handleSubmit}>` + 필드 + 제출 |
| `renderActionButton` | `render_action_button.go:21` | `<button onClick={() => mut.mutate({})}>` |
| `renderActionChildNodes` | `render_action_child_nodes.go:19` | action 폼 컨텍스트 children |
| `renderFilterUI` | `render_filter_ui.go:17` | 필터 컨트롤 |
| `renderSortUI` | `render_sort_ui.go:22` | 정렬 토글 |
| `renderPaginationUI` | `render_pagination_ui.go:17` | 페이지 이전/다음 |
| `renderComponentJSX` | `render_component_jsx.go:18` | `<CustomComponent data={...} />` |
| `renderStaticJSX` | `render_static_jsx.go:33` | 정적 DOM 트리 |
| `renderStaticActionJSX` | `render_static_action_jsx.go:30` | action 폼 내 정적 요소 |
| `renderFieldJSX` | `render_field_jsx.go` | `<input {...form.register(name)}>` |
| `renderBindJSX` | `render_bind_jsx.go:12` | data-bind |
| `renderEachJSX` | `render_each_jsx.go:13` | data-each |

## 임포트 수집

`importSet` (`import_set.go:6-16`)

```go
type importSet struct {
    react, useQuery, useMutation, useQueryClient bool
    useParams, useForm, useState                 bool
    components []string   // 참조 컴포넌트 목록
    customFile string     // "page-name.custom" 또는 ""
}
```

| 함수 | 파일 | 기능 |
|------|------|------|
| `collectImports` | `collect_imports.go:13-46` | 통합 분석 |
| `collectFetchImports` | `collect_fetch_imports.go:27` | Fetch 기반 요구 |
| `collectActionImports` | `collect_action_imports.go:11-23` | Action 기반 요구 |
| `renderImports` | `render_imports.go:60` | import 블록 문자열 |

`Dependencies(pages)` (`react_target_dependencies.go:7-22`) — 각 페이지 importSet을 돌아 npm 패키지 맵 생성:
- useQuery/useMutation/useQueryClient → `@tanstack/react-query ^5`
- useForm → `react-hook-form ^7`
- useParams → `react-router-dom ^6`

## 파라미터 처리

| 함수 | 파일 | 기능 |
|------|------|------|
| `renderParamArgs` | `render_param_args.go:21` | `ParamBind` → `{ key: value }` 객체 |
| `renderInfraApiArgs` | `render_infra_api_args.go:34` | pagination/sort/filter 인자 병합 |
| `paramSourceExpr` | `param_source_expr.go:16` | `"route.id"` → `"id"` |
| `extractRouteParamNames` | `extract_route_param_names.go` | 라우트 파라미터명 목록 |
| `renderUseParams` | `render_use_params.go` | `const { id } = useParams()` |
| `collectAllParams` | `collect_all_params.go:20` | Fetch/Action/children 전체 param 재귀 수집 |
| `collectFetchParamBinds` | `collect_fetch_param_binds.go:13` | Fetch 한정 |

## 이름 변환

| 함수 | 파일 | 예 |
|------|------|---|
| `toComponentName` | `to_component_name.go:14` | `"my-page"` → `"MyPage"` |
| `toLowerFirst` | `to_lower_first.go:14` | `"ListReservations"` → `"listReservations"` |
| `toUpperFirst` | `to_upper_first.go:14` | `"name"` → `"Name"` |

## 옵션

`GenerateOptions` (`generate_options.go:6-9`)

```go
type GenerateOptions struct {
    APIImportPath string   // default: "@/lib/api"
    UseClient     bool     // default: true ('use client' 디렉티브)
}
```

`DefaultOptions` (`default_options.go:6-11`) / `mergeOpt` (`merge_opt.go:5-11`).
orchestrator 은 `APIImportPath: "../api"`, `UseClient: false` 로 오버라이드.

## 엣지 케이스

| 시나리오 | 동작 |
|---------|------|
| Fetch 없음 | useQuery 블록 생략 |
| Action 없음 | useMutation/useForm 생략 |
| Children 없음 | `renderPageJSXFallback` — Fetch/Action 직접 나열 |
| Filter/Sort 없음 | 해당 UI 생략 |
| `route.X` 없음 | `useParams` 생략 |
| `custom.ts` 없음 | `import * as custom` 생략 |
| 중첩 Fetch | 재귀 처리 |
| 빈 컴포넌트명 | `orDefault()`로 기본값 |

## 결과 예시

```tsx
import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { useParams } from 'react-router-dom'
import { api } from '../api'

export default function MyReservationsPage() {
  const { id } = useParams()
  const queryClient = useQueryClient()
  const [page, setPage] = useState(1)

  const { data: listReservationsData, isLoading, error } = useQuery({
    queryKey: ['ListReservations', id, page],
    queryFn: () => api.ListReservations({ reservationId: id, page, limit: 20 }),
  })

  const createReservationForm = useForm()
  const createReservationMutation = useMutation({
    mutationFn: (data: any) => api.CreateReservation({ ...data, id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['ListReservations'] })
    },
  })

  return (
    <div>
      {isLoading && <div>로딩 중...</div>}
      {error && <div>오류가 발생했습니다</div>}
      {listReservationsData && (
        <section>
          {listReservationsData.items?.map((item: any, index: number) => (
            <div key={index}>{/* bind/each/static 렌더 */}</div>
          ))}
        </section>
      )}
      <form onSubmit={createReservationForm.handleSubmit((data) =>
        createReservationMutation.mutate(data))}>
        <input {...createReservationForm.register('title')} />
        <button type="submit">제출</button>
      </form>
    </div>
  )
}
```

## 파일 맵 (카테고리별)

### 공개 API
`generate.go`, `generate_page.go`, `generate_with.go`

### Target
`target_type.go`, `react_target_type.go`, `react_target_generate_page.go`, `react_target_file_extension.go`, `react_target_dependencies.go`, `default_target.go`

### 임포트
`import_set.go`, `collect_imports.go`, `collect_fetch_imports.go`, `collect_action_imports.go`, `render_imports.go`

### 훅 (Hook)
`render_page_hooks.go`, `render_fetch_hooks.go`, `render_use_query.go`, `render_page_mutations.go`, `render_use_mutation.go`, `render_form_hook.go`, `render_use_params.go`

### Fetch JSX
`render_fetch_jsx.go`, `render_fetch_jsx_body.go`, `render_fetch_jsx_flat_children.go`, `collect_fetch_ops.go`

### Action / Form JSX
`render_action_jsx.go`, `render_action_form.go`, `render_action_button.go`, `render_action_child_nodes.go`, `render_field_jsx.go`

### JSX 구조
`render_page_jsx.go`, `render_page_jsx_with_children.go`, `render_page_jsx_fallback.go`, `render_child_nodes.go`, `find_root_element.go`

### 데이터 바인딩
`render_bind_jsx.go`, `render_each_jsx.go`, `render_state_jsx.go`, `render_component_jsx.go`, `render_static_jsx.go`, `render_static_action_jsx.go`

### 인프라 UI
`render_filter_ui.go`, `render_sort_ui.go`, `render_pagination_ui.go`

### 파라미터
`render_param_args.go`, `render_infra_api_args.go`, `param_source_expr.go`, `extract_route_param_names.go`, `collect_all_params.go`, `collect_fetch_param_binds.go`

### 이름 변환 / 유틸
`to_component_name.go`, `to_lower_first.go`, `to_upper_first.go`, `cls_attr.go`, `indent_str.go`, `or_default.go`

### 옵션
`generate_options.go`, `default_options.go`, `merge_opt.go`

### 수집 / 기타
`collect_all_actions.go`, `deduplicate_actions.go`

## 설계 메모

- **문자열 빌더 기반** — AST/fmt 사용 안 함. `strings.Builder`로 TSX 직접 조립.
- **Target 확장 여지** — 인터페이스는 있으나 실제 구현은 React만. Vue/Svelte 추가하려면 `VueTarget` 같은 구현체만 붙이면 됨.
- **의존성 자동화** — importSet 분석이 `Dependencies()` 결과와 대칭이라 누락 없음.
- **재귀 분기 dispatcher** — `render_child_nodes.go` 의 switch가 모든 DSL 노드 렌더링의 중앙 허브.
- **invalidation 자동 배선** — action onSuccess가 모든 페이지 Fetch를 자동 무효화. 명시적 관계 선언 불필요.
