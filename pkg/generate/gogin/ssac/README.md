# internal/ssac/generator

SSaC 함수 정의를 **Go 서비스 핸들러 + 모델 인터페이스 + Handler 구조체** 로 변환하는 생성기.
Target 추상화를 통해 언어/프레임워크 교체 여지를 둠 (현재는 `GoTarget` 하나).

## 용어

**feature** = SSaC 파일의 상위 서브폴더(`specs/service/<feature>/*.ssac`). 내부 코드에서는 `Domain`.

## gogin과의 분업

| 책임 | 이 패키지 | `internal/gen/gogin` |
|------|-----------|---------------------|
| SSaC 시퀀스 → 핸들러 함수 body | **O** | — |
| DDL 스키마 → 모델 인터페이스 (`models_gen.go`) | **O** | — |
| feature별 Handler 구조체 (`handler.go`) | **O** | — |
| 서비스 파일 변환(func → method, receiver 주입) | — | **O** |
| Server/중앙 라우팅/main.go | — | **O** |
| 모델 구현(DB 접근) | — | **O** |
| Auth/middleware/authz | — | **O** |

한마디로:
- 이 패키지가 **인터페이스 + 핸들러 함수** 를 먼저 낳고,
- gogin이 그 위에 **Server·라우팅·모델 구현·초기화** 를 얹는다.

gogin의 `parseModelsGen()` 는 이 패키지가 만든 `models_gen.go` 를 읽어 구현체를 채운다.

## 진입점

| 함수 | 파일 | 설명 |
|------|------|------|
| `Generate(funcs, outDir, st, funcSpecs)` | `generate.go:12` | 기본 Target(`GoTarget`) 사용 |
| `GenerateWith(target, funcs, outDir, st)` | `generate_with.go:13` | Target 명시 |
| `GenerateModelInterfaces(funcs, st, modelDir)` | `go_target_generate_model_interfaces.go:15` | 모델 인터페이스만 별도 생성 |
| `GenerateHandlerStruct(funcs, st, outDir)` | `go_target_generate_handler_struct.go:17` | feature별 `handler.go` 생성 |

orchestrator 경로: `orchestrator/gen_ssac.go:60, 82` 에서 순차 호출.

## 산출물

| 대상 | 경로 | 파일 | 함수 |
|------|------|------|------|
| HTTP 핸들러 | `{outDir}/{feature}/{name}.go` | 개별 함수별 | `generateAndWrite:15` |
| Subscribe 핸들러 | `{outDir}/{feature}/{name}.go` | 개별 함수별 | `generateAndWrite:15` |
| 모델 인터페이스 | `{outDir}/model/models_gen.go` | 1개 통합 | `GenerateModelInterfaces:16` |
| Handler 구조체 | `{outDir}/{feature}/handler.go` | feature별 1개 | `GenerateHandlerStruct:18` |

feature 가 비어있는 함수는 `{outDir}/{name}.go` (기본 feature "service" 취급).

## Target 추상화

`target_interface.go:11`

```go
type Target interface {
    GenerateFunc(sf ServiceFunc, st *SymbolTable) ([]byte, error)
    GenerateModelInterfaces(funcs []ServiceFunc, st *SymbolTable, outDir string) error
    GenerateHandlerStruct(funcs []ServiceFunc, st *SymbolTable, outDir string) error
    FileExtension() string
}
```

구현: `GoTarget` (`go_target_type.go:8`)

```go
type GoTarget struct {
    FuncSpecs []funcspec.FuncSpec  // @error 에러 상태 코드 매핑용
}
```

| 메서드 | 파일 |
|--------|------|
| `GenerateFunc` | `go_target_generate_func.go:11` (HTTP/Subscribe 분기) |
| `GenerateModelInterfaces` | `go_target_generate_model_interfaces.go:15` |
| `GenerateHandlerStruct` | `go_target_generate_handler_struct.go:17` |
| `FileExtension` | `go_target_file_extension.go:6` → `".go"` |

## 핵심 흐름

### HTTP 함수 생성

```
GoTarget.GenerateFunc
  └─ sf.Subscribe == nil 분기
      └─ generateHTTPFunc                   (go_target_generate_http_func.go:11)
          ├─ analyzeHTTPFunc                (analyze_http_func.go:10)
          │    → httpFuncContext {pathParams, requestParams, flags…}
          ├─ buildHTTPFuncBody              (build_http_func_body.go:13)
          │    ├─ writeHTTPSequences        (write_http_sequences.go:13)
          │    │    └─ for each seq:
          │    │         └─ buildTemplateData + 템플릿 실행
          │    └─ 응답/에러 반환 코드
          ├─ collectImports                 (collect_imports.go:10)
          ├─ filterUsedImports              (filter_used_imports.go:16)
          └─ assembleGoSource               (assemble_go_source.go:11)
               → package + imports + body → gofmt
```

### Subscribe 함수 생성

```
sf.Subscribe != nil
  └─ generateSubscribeFunc                   (go_target_generate_subscribe_func.go:13)
      ├─ renderStructDefs                   (메시지 구조체)
      ├─ buildSubscribeFuncBody             (build_subscribe_func_body.go:13)
      │    └─ writeSubscribeSequences       (write_subscribe_sequences.go:13)
      ├─ collectSubscribeImports            (collect_subscribe_imports.go:12)
      └─ assembleGoSource
```

### 모델 인터페이스 생성

```
GenerateModelInterfaces
  ├─ collectModelUsages                     (collect_model_usages.go:7)
  │    └─ 각 함수의 seq.Model.Method 호출 수집
  ├─ deriveInterfaces                       (derive_interfaces.go:7)
  │    └─ 모델별로 필요한 메서드 파생
  │        ├─ deriveMethod                  (derive_method.go:10)
  │        └─ deriveReturnType              (derive_return_type.go:11)
  │             · wrapper("Page"/"Cursor") → pagination.X[T]
  │             · cardinality=="many" + hasQueryOpts → ([]T, int, error)
  │             · cardinality=="one" → (*T, error)
  ├─ renderInterfaces                       (render_interfaces.go:7)
  │    └─ renderSingleInterface             (render_single_interface.go:10)
  │         · WithTx(*sql.Tx) XxxModel
  │         · 각 파생 메서드
  └─ os.WriteFile({outDir}/model/models_gen.go)
```

### Handler 구조체 생성

```
GenerateHandlerStruct
  ├─ collectDomainModels                    (collect_domain_models.go:8)
  │    └─ feature별로 사용 모델 집합
  ├─ sortDomainModels                       (sort_domain_models.go:7)
  └─ writeHandlerFile                       (write_handler_file.go:21)
       · type Handler struct {
       ·   DB *sql.DB
       ·   CourseModel model.CourseModel
       ·   …
       · }
```

## 시퀀스 타입 → 템플릿

`go_templates.go:5` 에 정의된 템플릿 ID:

| 시퀀스 타입 | HTTP 템플릿 | Subscribe 템플릿 | 용도 |
|-----------|-----------|----------------|------|
| `SeqGet` | `get` | `sub_get` | 조회 |
| `SeqPost` | `post` | `sub_post` | 삽입 |
| `SeqPut` | `put` | `sub_put` | 수정 |
| `SeqDelete` | `delete` | `sub_delete` | 삭제 |
| `SeqEmpty` | `empty` | `sub_empty` | null 가드 |
| `SeqExists` | `exists` | `sub_exists` | 존재 가드 |
| `SeqState` | `state` | `sub_state` | 상태 전이 |
| `SeqAuth` | `auth` | `sub_auth` | 인가 체크 |
| `SeqCall` | `call_with_result` / `call_no_result` | 동일 | @call 함수 |
| `SeqPublish` | `publish` | `sub_publish` | 큐 발행 |
| `SeqResponse` | `response` / `response_direct` | — | HTTP 응답 |

분기:
- HTTP — `template_name.go`
- Subscribe — `subscribe_template_name.go`

## Pagination 처리

`deriveReturnType` `derive_return_type.go:11`:

```
if usage.Result.Wrapper != "":   // "Page" or "Cursor"
    return "(*pagination.Page[T], error)" | "(*pagination.Cursor[T], error)"
elif cardinality == "many" && hasQueryOpts:
    return "([]T, int, error)"           // int = total
elif cardinality == "one":
    return "(*T, error)"
```

`buildHasTotal` `build_has_total.go:11`:
- `hasQueryInput(seq.Inputs) && result is slice && wrapper == ""` → `HasTotal = true`

관련 import:
- `needs_pagination_import.go:5` — 인터페이스에 `pagination.` 있으면 `models_gen.go`에 해당 import 추가.

## 주요 데이터 구조

| 구조체 | 파일 | 용도 |
|--------|------|------|
| `modelUsage` | `model_usage.go:7` | 함수 × 모델 × 메서드 호출 정보 |
| `derivedInterface` | `derived_interface.go:5` | 파생된 인터페이스 |
| `derivedMethod` | `derived_method.go:5` | 파생된 메서드 (이름/파라미터/반환) |
| `derivedParam` | `derived_param.go:5` | 파라미터 (name, goType) |
| `templateData` | `template_data.go:7` | 템플릿 실행 입력 |
| `httpFuncContext` | `http_func_context.go:7` | HTTP 함수 분석 컨텍스트 |
| `typedRequestParam` | `typed_request_param.go:5` | 요청 파라미터 (name, goType, extractCode) |
| `rawParam` | `raw_param.go:5` | 미가공 파라미터 |

## 파일 맵 (카테고리별 주요)

### HTTP 핸들러
`go_target_generate_http_func.go`, `analyze_http_func.go`, `build_http_func_body.go`, `write_http_sequences.go`, `build_template_data.go`, `collect_imports.go`, `filter_used_imports.go`

### 모델 인터페이스
`go_target_generate_model_interfaces.go`, `collect_model_usages.go`, `collect_model_usages_from_func.go`, `collect_models_for_func.go`, `derive_interfaces.go`, `derive_interface_for_model.go`, `derive_method.go`, `derive_return_type.go`, `render_interfaces.go`, `render_single_interface.go`, `needs_pagination_import.go`

### Subscribe 핸들러
`go_target_generate_subscribe_func.go`, `build_subscribe_func_body.go`, `write_subscribe_sequences.go`, `collect_subscribe_imports.go`

### Handler 구조체
`go_target_generate_handler_struct.go`, `collect_domain_models.go`, `sort_domain_models.go`, `write_handler_file.go`

### 공통 유틸
`generate_and_write.go`, `generate_with.go`, `assemble_go_source.go`, `field_type_resolver.go`

### 템플릿
`go_templates.go`, `template_data.go`, `template_name.go`, `subscribe_template_name.go`

### 테스트
`test_generate_model_interface_*_test.go`, `test_generate_page_*_test.go`, `test_generate_sort_allowlist_test.go`, `test_generate_full_example_test.go`, `test_generate_domain_package_test.go` 등

## 설계 메모

- **Target 추상화** — 이론적으로 언어 교체 가능. 현재 Go 구현만 있음.
- **템플릿 + 조립 혼합** — 시퀀스 단위는 템플릿(`go_templates.go`), 상위 구조는 문자열 조립 + `gofmt`.
- **선형 파이프라인** — 분석 → 파생 → 렌더 → 조립. 순환 의존 없음 (DAG).
- **인터페이스/구현 분리** — 모델은 이 패키지가 인터페이스만 만들고, gogin이 구현을 채움. DIP 구조.
- **파일 237개** — 평균 20~30줄. filefunc Q-룰(함수 body 10줄 이하)을 엄격히 적용한 결과로 극도로 잘게 분할됨.
