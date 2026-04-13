# Phase012 — OrchestratorInternalRemoval

> `internal/orchestrator/` 를 pkg 전용으로 전환. Phase004/008 의 이월 부채 청산.

## 목표

orchestrator 의 `internal/*` import 23건 전부 제거. gen 계열 배선을 `pkg/generate/*`, 타입은 `pkg/fullend`/`pkg/parser/*`/`pkg/contract` 로 치환. `parsed.go` 의 `SymbolTable` 필드 제거.

완료 후 `rg "internal/(gen|ssac/generator|stml/generator|genapi|contract)" internal/orchestrator/` 가 0 hit 이어야 함.

---

## 범위

### A. internal/ssac/generator (3 파일)

| orchestrator 파일 | 현재 참조 | 교체 대상 |
|------------------|---------|----------|
| `gen_ssac.go` | `ssacgenerator.Generate(...)` | `pkg/generate/ssac.Generate(...)` |
| `default_profile.go` | `ssacgenerator.Profile` 타입 | `pkg/generate/ssac.Profile` (필요 시 타입 이관) |
| `target_profile_model.go` | 동일 | 동일 |

### B. internal/stml/generator (3 파일)

| orchestrator 파일 | 교체 대상 |
|------------------|----------|
| `gen_stml.go` | `pkg/generate/stml.Generate` |
| `default_profile.go` | `pkg/generate/stml.Profile` |
| `target_profile_model.go` | 동일 |

### C. internal/genapi (8 파일)

`internal/genapi` 는 파싱·공통 타입(ServiceFunc, ParsedResult 등). `pkg/fullend.Fullstack` 또는 `pkg/parser/ssac.ServiceFunc` 로 이관.

| orchestrator 파일 | 참조 타입 | 교체 |
|------------------|---------|------|
| `parsed.go` | `genapi.ParsedResult` 필드 (특히 `SymbolTable`) | `pkg/fullend.Fullstack` 필드 참조 + `SymbolTable` 필드 제거 |
| `validate_with.go` | `genapi.*` | pkg 버전 |
| `run_cross_validate.go` | 동일 | 동일 |
| `append_ssac_after_ddl.go` | 동일 | 동일 |
| `inject_func_err_status_from_parsed.go` | 동일 | 동일 |
| (기타 4 파일) | 동일 | 동일 |

### D. internal/gen/gogin (2 파일)

| orchestrator 파일 | 교체 |
|------------------|------|
| `gen_authz.go` | `pkg/generate/gogin.GenerateAuthz` (또는 기존 pkg 함수) |
| `gen_state_machines.go` | `pkg/generate/gogin.GenerateStateMachines` |

### E. internal/contract (4 파일)

`pkg/contract` 는 이미 존재. orchestrator 만 import 경로 교체:

| orchestrator 파일 | 변경 |
|------------------|------|
| `gen_with.go` | `internal/contract` → `pkg/contract` |
| `trace_artifacts.go` | 동일 |
| `run_contract_validate.go` | 동일 |
| `restore_preserved.go` | 동일 |

---

## 작업 순서

### Step 1. 사전 분석

- `filefunc chain` 으로 각 orchestrator 파일이 참조하는 internal 심볼 전수 파악
- pkg 대응 심볼 존재 확인. 누락 심볼 있으면 **Step 2 이전에 pkg 에 이식**
- `default_profile.go` / `target_profile_model.go` 의 Profile 타입 위치 결정 (ssac/stml 생성기 패키지 내부 or 공용 pkg)

### Step 2. pkg 누락분 이식

Phase004 에서 복사되지 않은 심볼을 pkg 로 이식. 최소한:
- `pkg/generate/ssac.Profile`
- `pkg/generate/stml.Profile`
- `pkg/generate/gogin.GenerateAuthz` (없으면 신설)
- `pkg/generate/gogin.GenerateStateMachines` (없으면 신설)

### Step 3. internal/contract → pkg/contract 4 파일

가장 단순 (import 경로 교체). 먼저 수행해 baseline 확인.

### Step 4. internal/gen/gogin 2 파일

gen_authz.go, gen_state_machines.go 교체.

### Step 5. internal/genapi 8 파일

**가장 복잡**. `ParsedResult` 가 orchestrator 전반에서 쓰이므로, `parsed.go` 의 `SymbolTable` 필드 제거와 맞물림. 다음 순서:
1. `parsed.go` 에서 `SymbolTable` 필드 삭제 + pkg/fullend 타입으로 교체
2. 각 사용처(validate_with, run_cross_validate 등) 수정
3. 빌드 통과 확인

### Step 6. internal/ssac/generator + internal/stml/generator (6 파일)

gen_ssac.go / gen_stml.go / default_profile.go / target_profile_model.go 2 쌍 교체.

### Step 7. 검증

- `rg "internal/(gen|ssac/generator|stml/generator|genapi|contract)" internal/orchestrator/` → 0 hit
- `go build ./pkg/... ./internal/... ./cmd/...`
- `go vet ./...`
- `fullend validate dummys/gigbridge/specs` 성공
- `fullend gen dummys/gigbridge/specs /tmp/x` 성공
- `scripts/structural_metrics.go` 재실행해 지표 비교

---

## 주의사항

### R1. Phase014 선행 아님

이 Phase 는 **orchestrator 만** 대상. `internal/gen/*`, `internal/ssac/generator`, `internal/stml/generator` 자체 삭제는 Phase014. 본 Phase 는 "의존 끊기" 까지.

### R2. pkg 누락분 이식 시 범위 한정

Step 2 에서 "그냥 복붙" 이 아니라 Phase010 의 Decide* 수렴 원칙 유지. 새 심볼도 2-depth 이내 if-else / Decide* 로.

### R3. Profile 타입 이관 위치 고민

`default_profile.go` / `target_profile_model.go` 의 Profile 이 ssac/stml generator 내부에 있을지, orchestrator 쪽에 둘지 결정. 후보:
- A안: `pkg/generate/ssac.Profile` + `pkg/generate/stml.Profile` (generator 소유)
- B안: `internal/orchestrator/profile.go` (orchestrator 소유)
- 권장: A안 — 생성기가 자신의 Profile 정의

### R4. parsed.go SymbolTable 제거는 breaking

`ParsedResult.SymbolTable` 필드 제거 시 외부 호출자(있을 경우) 영향. 사전에 `rg "\.SymbolTable"` 전수 조사 후 진행.

---

## 완료 조건 (Definition of Done)

- [x] `internal/orchestrator/` 내 `internal/{gen,ssac/generator,stml/generator,genapi,contract}` import 0건
- [~] `parsed.go` 에서 `SymbolTable` 필드 삭제 — **유보**: validate 측이 SymbolTable 사용. `ParsedSSOTs` 를 orchestrator 로 인라인 (genapi 의존만 제거). 진정한 SymbolTable 제거는 validate 측 마이그레이션 필요 (Phase013/014 후속)
- [x] `go build ./pkg/... ./internal/... ./cmd/...` 통과
- [x] `go vet ./...` 통과
- [x] `go test ./pkg/...` 통과
- [x] `fullend gen dummys/{gigbridge,zenflow}/specs` 성공 (10 artifacts)
- [ ] `scripts/structural_metrics.go` 재실행 — Phase013 또는 014 와 함께
- [ ] 커밋: `refactor(orchestrator): internal/* 의존 전수 제거 (pkg 전환 완결)`

### 부산물

- `pkg/generate/react.Generate` stub 해소 (gen_glue 가 pkg/generate.Generate 호출하면서 활성화 필요)
- `pkg/generate/gogin/ssac/*.go` 의 `internal/funcspec` import 5건을 `pkg/parser/funcspec` 으로 교체
- `internal/orchestrator/inject_func_err_status_from_parsed.go` 삭제 (dead code: ground.Build 가 동등 기능 수행)
- `internal/orchestrator/parsed_ssots.go` 신설 (genapi.ParsedSSOTs 인라인)
- `internal/orchestrator/determine_pkg_module_path.go` 신설 (manifest.ProjectConfig 기반)

## 의존

- Phase011 완료 ✅ (pkg 안정성 기준선 확보)

## 다음 Phase

- **Phase013** — 생성 품질 (zenflow type-mismatch 해소)
- **Phase014** — internal/* 일괄 삭제 (본 Phase 가 의존 끊으면 가능)
