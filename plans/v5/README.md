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
| 012 | OrchestratorInternalRemoval | ✅ | orchestrator `internal/{gen,ssac/generator,stml/generator,genapi,contract}` import 0 달성. 커밋 `403eae5` |
| 013 | GenerateQualityFixing | ✅ | zenflow spec 정합화로 14 type-mismatch 해소. `FULLEND_LOCAL_PATH` env replace 옵션. 커밋 `8ed9352` |
| **014** | **InternalGenRemoval** | **대기** | `internal/gen/*`, `ssac/generator`, `stml/generator`, `genapi`, `contract` 일괄 삭제 (pkg/cmd/crosscheck 선행 정리 필요) |
| 015 | TemplateAndResiduals | 대기 | `*_template.go` → `text/template` + pkg→internal 역의존(7) 제거 + filefunc F1/F2 정책 결정 |
| **016** | **CrosscheckStrengthening** | **대기** | 정합성 규칙 6종 추가 — validate/crosscheck 먼저 강화해 spec 결함을 사전 검출 (DDL CHECK↔INSERT / DEFAULT FK↔seed / claims↔DDL / SSaC role↔OPA / @empty↔nilable / @call 인자 타입) |
| 017 | RuntimeBugFixing | 대기 | 실측 런타임 버그 4종 수정 (zenflow 403, OPA path, DSN, gigbridge seed). Phase016 validate 결과를 작업 리스트로 활용 |
| 018 | DDLPipelineIntegration | 대기 | DDL 위상정렬 + `schema.sql` 통합 산출 + `DEFAULT N FK` auto nobody seed |

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

**Phase016 — CrosscheckStrengthening** (권장 — validate/crosscheck 먼저 강화) 또는 **Phase014 — InternalGenRemoval** (구조 정리 트랙)

### Phase016 → 017 → 018 (실측 버그 트랙 — 검증 먼저)
2026-04-14 docker+hurl 실측에서 드러난 결함 대응을 **3개 Phase 로 분할**. **"도구 먼저 고치고 그 도구로 버그 잡기"** 원칙:
- **016 CrosscheckStrengthening** (중): 정합성 규칙 6종 추가 — spec 결함을 validate 단계에서 사전 검출
- **017 RuntimeBugFixing** (소-중): zenflow 403, OPA path, DSN 기본값, gigbridge seed — 016 의 validate 결과 + validate 로 못 잡는 순수 런타임 버그
- **018 DDLPipelineIntegration** (중): DDL 위상정렬, `schema.sql` 통합, auto nobody seed — 기능 추가

권장 순서: **016 → 017 → 018**. 016/018 은 순서 바꿔도 되나 016 먼저 하면 017 조사 범위가 줄어듦.

### Phase014 (구조 정리 트랙)
Phase012 후 잔여 cmd/crosscheck/reporter 의 internal 의존 정리 후 일괄 삭제. Phase015 로 이어짐.

상세: `plans/v5/Phase{016,017,018,014,015}-*.md`.

---

## 구조 건전성 실측 (2026-04-14)

Phase013 완료 후 실측 결과 — 전 지표 개선:

| 지표 | internal/gen | pkg/generate | 변화 |
|------|-------------|--------------|------|
| 평균 복잡도 (cyclomatic) | 5.28 | **3.95** | **-25%** ✅ |
| 고복잡(16+) 함수 비율 | 4.3% | **1.6%** | **-63%** ✅ |
| 단순(1) 함수 비율 | 13.8% | **23.5%** | **+70%** ✅ |
| 평균 매개변수 | 2.56 | **2.23** | -13% ✅ |
| `*WithDomains` 중복 | 4 | **0** | -100% ✅ |
| `Decide*` 수렴점 | 0 | **3** | Phase010 ✅ |
| orchestrator 대상 internal 의존 | 23 | **0** | Phase012 ✅ |

상세: `reports/structural-evaluation-2026-04-14.md`.

기능 결함 (OPA owners 정확 원인, crosscheck 사각지대) 은 구조와 별개 — Phase016 대상.

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
