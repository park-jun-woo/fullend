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
| 004 | InternalToPkgMigration | ✅ | pkg 이식 + orchestrator 잔여 의존 제거 (Phase012 에서 완결). 커밋 `403eae5` |
| 005 | GroundConvergence | ✅ | ssac generator 의 `validator.SymbolTable` → `rule.Ground` 전수 치환 |
| 006 | FlatRemovalAndGoginActivation | ✅ | Flat mode 제거 + gogin.Generate stub 해소 |
| 007 | ReactHurlActivation | ✅ | react.Generate + hurl.Generate stub 해소 |
| 008 | OrchestratorWiring | ✅ | gen_glue + genSSaC/genSTML 까지 완결. 커밋 `403eae5` |
| 009 | StructuralCleanup | ✅ | generateMain struct 수렴. Phase010 Decide\* 수렴으로 마감. 커밋 `aab3c48` |
| 010 | ToulminPoints → **Decide\* 순수 함수 수렴** | ✅ | 2-depth 기준점 적용 결과 Toulmin 미채택. 커밋 `aab3c48` |
| 011 | DummyRegressionValidation | ✅ | Tier 1 통과, gigbridge 빌드 ✓, zenflow 14 type-mismatch (Phase013 이월). 커밋 `304769d` (v0.1.9) |
| 012 | OrchestratorInternalRemoval | ✅ | orchestrator `internal/{gen,ssac/generator,stml/generator,genapi,contract}` import 0 달성. 커밋 `403eae5` |
| 013 | GenerateQualityFixing | ✅ | zenflow spec 정합화로 14 type-mismatch 해소. `FULLEND_LOCAL_PATH` env replace. 커밋 `8ed9352` |
| 014 | InternalGenRemoval | ⛔ 실행 금지 | 사용자 지시로 **현재 진행 안함** — internal/* 삭제는 v5 이후 별도 판단 |
| 015 | TemplateAndResiduals | ✅ | Part A template 전환 (main/query_opts → embed + text/template, bit-level 동일), Part B pkg→internal 의존 **0 hit** (internal/policy→pkg/parser/rego 치환, adapt_policies 삭제), Part C docs/filefunc-policy.md 확정. baseline 37→36 |
| 016 | CrosscheckStrengthening | ✅ | 정합성 규칙 6종 (X-74~X-79) 추가. 파서 확장 (DDL Defaults/Seeds, FuncSpec.ResponsePointer). 커밋 `e870181` (v0.1.10) |
| 017 | RuntimeBugFixing | ✅ | 런타임 버그 4종 (zenflow 403 ← OPA claims `org_id` 누락 / OPA path 디렉토리+fallback / DSN DATABASE_URL+모듈명 / gigbridge seed). gigbridge smoke 12/12. **미커밋** |
| 018 | DDLPipelineIntegration | ✅ | DDL 위상정렬 (Kahn) + `schema.sql` 통합 산출 + `DEFAULT N FK` auto nobody seed (opt-in). gigbridge schema.sql 1회 실행으로 DB 초기화. **미커밋** |

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

**v5 로드맵 완료** (Phase001~013, 015~018). Phase014 만 ⛔ 실행 금지 (사용자 지시).

### 후속 판단 필요 항목 (v6 후보)

- **Phase014 internal/\* 일괄 삭제** — 재활성화 시점은 별도 판단
- **zenflow ListWorkflows 500**: SSaC `@get` 의 org 스코프 자동 주입 설계 (Phase017 에서 발견)
- **generator 비결정성**: 같은 입력에 대해 map 순회 기반으로 출력 순서가 변동. deterministic codegen 이슈
- **filefunc 정책**: `docs/filefunc-policy.md` 제출 + filefunc 리포에 `//ff:group` 이슈

### Phase016 → 017 → 018 (실측 버그 트랙 — 마감 ✅)
- **016 ✅ CrosscheckStrengthening**: 정합성 규칙 6종 (X-74~X-79) 추가. 파서 확장. 커밋 `e870181` / v0.1.10
- **017 ✅ RuntimeBugFixing**: zenflow ActivateWorkflow 403 root cause (OPA claims.org_id 누락) + OPA path fallback + DSN DATABASE_URL + gigbridge seed. gigbridge smoke 12/12.
- **018 ✅ DDLPipelineIntegration**: DDL 위상정렬 + schema.sql 통합 + auto nobody seed (opt-in).

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
