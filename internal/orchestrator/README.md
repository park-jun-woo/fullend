# internal/orchestrator

CLI 명령(`validate`, `gen`, `status`, `chain`)의 **최상위 실행 제어**와 **SSOT 1회 파싱 · N회 공유**를 담당.
개별 파서·검증기·생성기를 연결하는 접착층.

## 용어 주석

**feature** = `specs/service/<폴더>/*.ssac` 서브폴더명. 내부 식별자로는 `Domain`.

## 진입점

| CLI 명령 | 공개 함수 | 파일 | 설명 |
|---------|---------|------|------|
| `validate` | `Validate(root, detected, skipKinds)` | `validate.go:15` | SSOT 검증 |
| `gen` | `Gen(specsDir, artifactsDir, skipKinds, reset)` | `gen_cmd.go:12` | validate 선행 후 코드 생성 |
| `status` | `Status(root, detected)` | `status_cmd.go:12` | SSOT 현황 집계 |
| `chain` | `Chain(specsDir, operationID)` | `chain_cmd.go:14` | operationID 연관 노드 추적 |

호출 구조:

```
CLI command (cmd/fullend/...)
     ↓
public func (Validate/Gen/Status/Chain)
     ↓
DetectSSOTs → ParseAll → [Validate|Generate|Aggregate|Trace]
     ↓
서브패키지 (pkg/validate, pkg/crosscheck, internal/gen, …)
```

## ParseAll — 1회 파싱 · N회 공유

**파일**: `parsed.go:22-100`
**반환**: `*genapi.ParsedSSOTs`

| SSOT 종류 | 파서 | 저장 필드 | 라인 |
|---------|------|----------|------|
| Config | `projectconfig.Load(root)` | `Config` | `parsed.go:31-35` |
| OpenAPI | `openapi3.NewLoader().LoadFromFile()` | `OpenAPIDoc` | `parsed.go:37-42` |
| DDL | `ssacvalidator.LoadSymbolTable(root)` | `SymbolTable` | `parsed.go:44-49` |
| SSaC | `ssacparser.ParseDir()` | `ServiceFuncs` | `parsed.go:51-56` |
| STML | `stmlparser.ParseDir()` | `STMLPages` | `parsed.go:58-63` |
| States | `statemachine.ParseDir()` | `StateDiagrams` + `StatesErr` | `parsed.go:65-72` |
| Policy | `policy.ParseDir()` | `Policies` | `parsed.go:74-79` |
| Func | `funcspec.ParseDir()` | `ProjectFuncSpecs` | `parsed.go:81-86` |
| FullendPkg | `funcspec.ParseDir(pkgRoot)` | `FullendPkgSpecs` | `parsed.go:89-93` |
| Model | (경로만 저장) | `ModelDir` | `parsed.go:95-97` |

`--skip`에 지정된 종류는 파싱 자체를 건너뜀. 파싱 실패 시 해당 필드는 nil, 에러는 명시적으로 기록하지 않음 (States는 예외적으로 `StatesErr`에 보존).

공유 패턴:
- `Validate` → `ParseAll` → `ValidateWith`
- `Gen` → `GenWith` → `ParseAll` → `ValidateWith` → `runCodegenSteps`
- `Status` → `ParseAll` → 집계 함수들
- `Chain` → `ParseAll` → trace 함수들

## validate 흐름

`validate.go:15-25` + `validate_with.go:11-78`

```
Validate(root, detected, skipKinds)
  ├─ ParseAll(root, detected, skip)
  └─ ValidateWith(root, detected, parsed, skip)
      ├─ for kind in allKinds:
      │    ├─ skip[kind] 시 reporter.Skip로 기록
      │    └─ kind별 validateXXX() 호출
      ├─ DDL 검증 뒤 SSaC 즉시 (appendSSaCAfterDDL.go:12)
      ├─ runCrossValidate(parsed)
      └─ runContractValidate(artifactsDir)
```

### 개별 검증 함수

| Kind | 함수 | 파일 | 라인 |
|------|-----|------|------|
| Config | `validateConfig` | `validate_config.go` | - |
| OpenAPI | `validateOpenAPI` | `validate_openapi.go:13` | 엔드포인트 수 + path 충돌 |
| DDL | `validateDDL` | `validate_ddl.go:14` | `pkg/validate/ddl` + sqlc 중복 검사 |
| SSaC | `validateSSaC` | `validate_ssac.go:16` | `pkg/validate/ssac` + `pkg/ground` |
| STML | `validateSTML` | `validate_stml.go:15` | `pkg/validate/stml` + `pkg/ground` |
| States | `validateStates` | `validate_states.go:14` | `pkg/validate/statemachine` |
| Policy | `validatePolicy` | `validate_policy.go:12` | 파일·규칙 개수 |
| Scenario | `validateScenarioHurl` | `validate_scenario.go:12` | `.feature` 차단 + `.hurl` 수집 |
| Func | `validateFunc` | `validate_funcspec.go:14` | `pkg/validate/funcspec` |
| Model | `validateModel` | `validate_model.go` | - |

### 교차 검증

`runCrossValidate()` `run_cross_validate.go:15-55`:
- **조건**: OpenAPI + DDL + SSaC 전부 존재.
- `pkg/crosscheck.Run()` 호출 — Toulmin 기반 교차 정합성.
- 에러/경고 카운트는 `countCrossErrors:52`.

### 계약 검증

`runContractValidate()` `run_contract_validate.go:15-72`:
- `artifacts/` 의 `//fullend:gen ownership=…` 디렉티브 순회.
- `broken`/`orphan` 발견 시 fail.
- `artifacts/` 자체가 없거나 디렉티브가 없으면 skip.

## gen 흐름

`gen_cmd.go:12` + `gen_with.go:17-90`

```
Gen(specsDir, artifactsDir, skipKinds, reset)
  └─ GenWith(DefaultProfile(), specsDir, artifactsDir, skipKinds, reset)
      ├─ DetectSSOTs(specsDir)
      ├─ ParseAll(…)
      ├─ ValidateWith(…)               ← validate 실패 시 즉시 반환
      ├─ snapshotPreserved()            ← reset=false 일 때 본문 백업
      ├─ runCodegenSteps(…)
      └─ restorePreserved(snapshot)     ← reset=false 일 때 복원
```

### runCodegenSteps 단계

`run_codegen_steps.go:15-54`

| 단계 | 함수 | 파일 | 조건 |
|------|-----|------|------|
| sqlc | `genSqlc` | `gen_sqlc.go:12` | DDL 존재 |
| oapi-gen | `genOpenAPI` | `gen_openapi.go:14` | OpenAPI 존재 |
| ssac-gen | `genSSaC` | `gen_ssac.go:16` | SSaC 존재 |
| ssac-model | (genSSaC 내부) | `gen_ssac.go:81-89` | SSaC 존재 |
| stml-gen | `genSTML` | `gen_stml.go:16` | STML 존재 |
| glue-gen | `genGlue` | `gen_glue.go:14` | 항상 — `internal/gen.Generate()` 호출 |
| hurl-gen | (체크만) | `run_codegen_steps.go:38-42` | `tests/smoke.hurl` 존재 시 |
| state-gen | `genStateMachines` | `gen_state_machines.go:14` | States 존재 |
| authz-gen | `genAuthz` | `gen_authz.go:14` | Policy 존재 |
| func-gen | `genFunc` | `gen_func.go:13` | Func 존재 |

### 최종 산출 호출

`genGlue` (`gen_glue.go:14-43`) 가 `GenConfig`와 `STMLGenOutput`을 조립한 뒤 `gen.Generate(parsed, cfg, stmlOut)` 실행 — 이후 흐름은 `internal/gen/README.md` 참조.

### Preserve 메커니즘

`--reset=false` (기본) 시:
- `snapshotPreserved()` — `//fullend:gen ownership=preserve` 마커가 있는 함수 본문을 메모리에 저장.
- `restorePreserved()` — 재생성 후 해당 body를 되돌림.
- contract 해시 변경 시 `[WARN]` 로 리포트에 기록.

## status 흐름

`status_cmd.go:12-85`

```
Status(root, detected)
  ├─ ParseAll(root, detected, nil)
  └─ StatusLine 배열 생성
```

각 SSOT 종류별 집계 함수:

| SSOT | 집계 함수 | 파일 |
|------|----------|------|
| OpenAPI | `countEndpoints(doc)` | `count_endpoints.go` |
| DDL | `countDDLColumns(tables)` | `count_ddl_columns.go` |
| SSaC | `len(parsed.ServiceFuncs)` | (인라인) |
| STML | `len(parsed.STMLPages)` | (인라인) |
| States | `countTransitions(diagrams)` | `count_transitions.go` |
| Policy | `countPolicyRules(policies)` | `count_policy_rules.go` |
| Scenario | `filepath.Glob("*.hurl")` | (인라인) |
| Func | `countFuncStubs(specs)` | `count_func_stubs.go` |

Model은 보조 취급이라 status에 표시하지 않음.

## chain 흐름

`chain_cmd.go:14-110`

```
Chain(specsDir, operationID)
  ├─ DetectSSOTs + ParseAll
  ├─ SSaC 함수 목록에서 operationID 매칭 → matchedFunc
  └─ 9개 trace 함수 호출 → ChainLink 배열
```

### Trace 함수 (`trace_*.go`)

| Trace | 파일 | 추적 대상 | 선택 기준 |
|-------|------|---------|---------|
| `traceOpenAPI` | `trace_openapi.go:13` | OpenAPI operation | operationID 매칭 |
| `traceSSaC` | (ssac parser) | SSaC 함수 | 함수명 매칭 |
| `traceDDL` | `trace_ddl.go:15` | DDL 테이블 | seq.Model → 테이블명 |
| `tracePolicy` | `trace_policy.go:14` | Rego 규칙 | @auth seq.Resource |
| `traceStates` | `trace_states.go:11` | State diagram | @state seq.DiagramID |
| `traceFuncSpecs` | `trace_func_specs.go:13` | FuncSpec | @call seq.Model |
| `traceHurlScenarios` | `trace_hurl_scenarios.go:12` | `.hurl` 파일 | endpoint path 매칭 |
| `traceSTML` | `trace_stml.go:12` | STML 페이지 | operationID 속성 |
| `traceArtifacts` | `trace_artifacts.go:16` | 생성 코드 | SSOT 경로 + 함수명 |

### ChainLink 구조

`chain_link.go`:
- `Kind`: OpenAPI / SSaC / DDL / Rego / StateDiag / FuncSpec / Hurl / STML / Handler / Model / Authz / States
- `File`: 상대 경로
- `Line`: 1-based (미상 시 0)
- `Summary`: 요약 문자열
- `Ownership`: `""`(SSOT) / `"gen"` / `"preserve"`

## 주요 파일 맵

### 공개 커맨드

- `validate.go` / `validate_with.go` — validate 본체
- `gen_cmd.go` / `gen_with.go` — gen 본체
- `status_cmd.go` — status 본체
- `chain_cmd.go` — chain 본체

### 감지 · 파싱

- `detect_ssots.go:12` — 루트에서 SSOT 종류 감지
- `detected_ssot.go` — `DetectedSSOT{Kind, Path}`
- `ssot_kind.go` — `SSOTKind` 상수
- `parsed.go:22` — `ParseAll()`

### 검증 스텝

- `validate_*.go` — 종류별 validator
- `run_cross_validate.go` — 교차
- `run_contract_validate.go` — 계약

### 코드젠 스텝

- `run_codegen_steps.go` — 10단계 오케스트레이션
- `gen_sqlc.go` / `gen_openapi.go` / `gen_ssac.go` / `gen_stml.go` / `gen_glue.go` / `gen_state_machines.go` / `gen_authz.go` / `gen_func.go`

### Chain trace

- `trace_*.go` × 9
- `build_state_chain_link.go:9`

### 보조 유틸

- `find_ssac_file.go:13` — feature 구조 vs 평면 경로 결정
- `scan_func_imports.go:14` — SSaC의 func 패키지 import 스캔
- `copy_func_package.go:17` — func/ → artifacts 파일 복사
- `restore_preserved.go:15` — 본문 복원 + 계약 경고
- `inject_func_err_status_from_parsed.go:15` — funcspec `@error` → symbol table 주입
- `append_ssac_after_ddl.go:12` — DDL 검증 후 SSaC 검증 즉시 추가
- `default_profile.go:12` — `DefaultProfile()` Go + React 반환

### 설정 결정

- `detect_db_engine.go`
- `determine_module_path.go`
- `find_fullend_pkg_root.go`
- `run_exec.go:34` — 외부 명령 실행 래퍼 (timeout + skip-if-missing)

### 타입

- `chain_link.go`, `detected_ssot.go`, `ssot_kind.go`, `target_profile_model.go`

### 집계

- `count_endpoints.go`, `count_ddl_columns.go`, `count_transitions.go`, `count_policy_rules.go`, `count_func_stubs.go`
- `find_func_spec_link.go:14`, `resolve_func_spec_path.go`, `grep_line.go`
- `missing_ssot_step.go:9` — 미감지 SSOT의 StepResult

## 리포트 · Exit 전략

`reporter.StepResult`:
- `Name`, `Status(Pass/Fail/Skip)`, `Summary`, `Errors[]`, `Suggestions[]`

에러 수집:
1. 개별 validator 각자 `Errors` 채움.
2. 교차 검증 — `pkgcross.Run`의 `ValidationError[]` 변환.
3. 계약 검증 — broken/orphan 문자열.
4. Preserve 복원 — `[WARN]` 접두로 경고 기록.

Exit:
- `report.HasFailure() == true` → `GenWith` 반환 `bool == false` → CLI 비정상 종료.
- Skip 상태는 실패로 보지 않음.

## 설계 요약

| 원칙 | 구현 위치 |
|------|----------|
| 1회 파싱 재사용 | `parsed.go` (`ParseAll`) |
| validate 필수 선행 | `gen_with.go:40` |
| `--skip` 일관 처리 | `validate_with.go:28-34` |
| operation 단위 체이닝 | `chain_cmd.go` + `trace_*.go` |
| 보존-재생성 사이클 | `snapshotPreserved` + `restorePreserved` |
| 외부 도구 호출 일관화 | `run_exec.go` (timeout + missing skip) |
