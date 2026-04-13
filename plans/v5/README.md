# plans/v5 — pkg 기반 리팩토링 로드맵

## 전체 목표

internal 코드젠을 프로덕션 레벨로 끌어올리는 기반 다지기.
**구조 건전성 > 기능 완성도** (이 단계 한정).
2-depth 이내는 if-else, 초과는 Toulmin.

## Phase 상태 (2026-04-13 시점)

| # | 제목 | 상태 | 비고 |
|---|------|------|------|
| 001 | FullendSplitAndFeatureRename | ✅ | `pkg/fullend/` 분리, `Domain→Feature` |
| 002 | GroundExpansion (+ GroundDesign) | ✅ | Ground 신 필드: Models/Tables/Ops/ReqSchemas. iface + sqlc 파서 신설 |
| 003 | ValidateCrosscheckAlignment | ✅ | 5개 규칙 Ground 마이그 (S-48, X-12, X-13, X-55, X-9/X-10). `populate_model_lookup` 제거 |
| 004 | InternalToPkgMigration | 🟡 | pkg/contract, pkg/generate/{gogin, react, hurl, ssac, stml} 이식 완료. orchestrator 잔여 의존은 Phase012 에서 해소 |
| 005 | GroundConvergence | ✅ | ssac generator 의 `validator.SymbolTable` → `rule.Ground` 전수 치환 |
| 006 | FlatRemovalAndGoginActivation | ✅ | Flat mode 제거 + gogin.Generate stub 해소 (MainGenInput struct 포함) |
| 007 | ReactHurlActivation | ✅ | react.Generate + hurl.Generate stub 해소. `adapt_policies.go` (rego→internal policy 어댑터) |
| 008 | OrchestratorWiring | 🟡 | gen_glue.go 배선 교체 완료. **Phase012 에서 genSSaC/genSTML 까지 완결** |
| 009 | StructuralCleanup | ✅ | generateMain struct 수렴. Toulmin 도입은 Phase010 에서 처리 완료 |
| 010 | ToulminPoints → **Decide\* 순수 함수 수렴** | ✅ | 2-depth 기준점 적용 결과 Toulmin 미채택. 3 포인트 모두 `Decide*(facts) Decision` 순수 함수로 수렴. 커밋 `aab3c48` |
| 011 | DummyRegressionValidation | ✅ | Tier 1 통과, gigbridge 빌드 ✓, zenflow 17 type-mismatch (Phase013 이월). `reports/metrics-phase011.md` 생성. 커밋 `304769d` (v0.1.9) |
| **012** | **OrchestratorInternalRemoval** | **대기** | orchestrator 내 internal/* import 23건 완전 제거 (genSSaC/genSTML/parsed.go/gen_authz 등) |
| 013 | GenerateQualityFixing | 대기 | zenflow 17 type-mismatch 해소 + dummy go.mod replace 자동 주입 |
| 014 | InternalGenRemoval | 대기 | `internal/gen/*`, `internal/ssac/generator`, `internal/stml/generator`, `internal/genapi`, `internal/contract` 일괄 삭제 |
| 015 | TemplateAndResiduals | 대기 | `main_template.go` 등 `text/template` 화 + 잔여 internal 의존 정리 + filefunc F1/F2 완화 논의 |

## 검증 상태 (Phase011 기준선)

- `go build ./pkg/... ./internal/... ./cmd/...` ✅
- `go vet` ✅
- `go test ./pkg/...` 전수 통과 ✅
- `fullend validate dummys/{gigbridge,zenflow}/specs` 정상 동작 ✅
- `fullend gen dummys/{gigbridge,zenflow}/specs dummys/*/artifacts` 성공 — **pkg/generate 경유** (glue-gen ✓) ✅
- `cd dummys/gigbridge/artifacts/backend && go build ./...` ✅ (replace 수동 추가 후)
- `cd dummys/zenflow/artifacts/backend && go build ./...` ✗ (17 type-mismatch, Phase013 대상)
- `filefunc validate`: 내 파일 위반 0, baseline 37 잔존 (기존 파일들)

## 핵심 아키텍처 결정

- **pkg/fullend.Fullstack** = SSOT 파싱 결과 canonical 컨테이너
- **pkg/rule.Ground** = validate + crosscheck + generate 공유 조회 계층
  - 기존 평탄: `Lookup/Types/Pairs/Schemas` (validate 주도)
  - 신규 구조: `Models/Tables/Ops/ReqSchemas` (generate 주도)
- **복사 방식** — internal 은 참조용으로 유지, Phase014 에서 일괄 삭제
- **기술 스택 단일** — `gogin/react` 디렉토리는 수요 기반 확장 슬롯 (nisabit.com 운영 데이터 기반)
- **결정 로직 수렴** — `Decide*(facts) Decision` 순수 함수 패턴 (Phase010). 2-depth 초과 시에만 Toulmin.

## 기준점 (사용자 확정, Phase010 에서 도입)

**"if-else 2-depth 이내에 해결 안되면 Toulmin"** — `depth 1` = flat chain, `depth 2` = 중첩 1회, `depth 3+` 부터 Toulmin 대상.
조건식 AND/OR 는 depth 미포함 (수평 확장). 본 기준은 이후 Phase 에도 계승.

## 다음 세션 시작 지점

**Phase012 — OrchestratorInternalRemoval**

`internal/orchestrator/` 내부 **23개 internal/* import** 전부 제거해 orchestrator 를 pkg 전용으로 전환.

대상 import 군:
- `internal/ssac/generator` — genSSaC, default_profile, target_profile_model (3 파일)
- `internal/stml/generator` — genSTML, default_profile, target_profile_model (3 파일)
- `internal/genapi` — parsed.go 타입 + validate_with/run_cross_validate/append_ssac_after_ddl/inject_func_err_status_from_parsed 등 (8 파일)
- `internal/gen/gogin` — gen_authz, gen_state_machines (2 파일)
- `internal/contract` — gen_with, trace_artifacts, run_contract_validate, restore_preserved (4 파일)

`parsed.go` 의 `SymbolTable` 필드 제거 포함.

상세: `plans/v5/Phase012-OrchestratorInternalRemoval.md`.

## 주요 산출 변경 (누적)

- `pkg/fullend/` — Fullstack 타입, ParseAll, DetectSSOTs, SSOTKind 등
- `pkg/parser/iface/`, `pkg/parser/sqlc/` 신설
- `pkg/rule/ground_types.go` — 11개 신 타입
- `pkg/ground/populate_{models,tables,ops,request_schemas}.go` — Ground 채움
- `pkg/contract/` — internal/contract 이식
- `pkg/generate/` — 전체 계층 (gogin/react/hurl + ssac/stml 서브). Phase010 Decide* 수렴 적용.
- `pkg/rule/model_ref_exists.go` — S-48 용 ModelRefExists warrant
- `scripts/structural_metrics.go` + `reports/metrics-phase011.md` — Phase011 지표

## 참고 문서

- `Phase002-GroundDesign.md` — Phase002 설계 리뷰 산출물 (SymbolTable 접근 패턴 전수 조사 + Ground 설계)
- `Phase010-{MethodDispatch,MainInit,ScenarioOrder}Design.md` — Decide* 3 포인트 설계 근거
- `reports/metrics-phase011.md` — 구조 건전성 비교 리포트 (internal vs pkg)
- `internal/gen/*/README.md` 9개 — internal 코드젠 분석 문서 (Phase001 이전에 작성, Phase014 에서 정리 대상)

## CLAUDE.md 규약 리마인더

- **.ffignore 수정 금지**
- **Co-Authored-By 금지**
- **버전 자동 bump** — `go install -ldflags "-X main.Version=v0.1.N"` (N 증가). 현재: v0.1.9.
- **커밋+푸시 시 민감 정보 확인 필수**
