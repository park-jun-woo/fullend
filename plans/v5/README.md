# plans/v5 — pkg 기반 리팩토링 로드맵

## 전체 목표

internal 코드젠을 프로덕션 레벨로 끌어올리는 기반 다지기.
**구조 건전성 > 기능 완성도** (이 단계 한정).
단순은 if-else, 복잡은 Toulmin.

## Phase 상태 (2026-04-13 시점)

| # | 제목 | 상태 | 비고 |
|---|------|------|------|
| 001 | FullendSplitAndFeatureRename | ✅ | `pkg/fullend/` 분리, `Domain→Feature` |
| 002 | GroundExpansion (+ GroundDesign) | ✅ | Ground 신 필드: Models/Tables/Ops/ReqSchemas. iface + sqlc 파서 신설 |
| 003 | ValidateCrosscheckAlignment | ✅ | 5개 규칙 Ground 마이그 (S-48, X-12, X-13, X-55, X-9/X-10). `populate_model_lookup` 제거 |
| 004 | InternalToPkgMigration | 🟡 | pkg/contract, pkg/generate/{gogin, react, hurl, ssac, stml} 이식 완료. stub 상태 시작 |
| 005 | GroundConvergence | ✅ | ssac generator 의 `validator.SymbolTable` → `rule.Ground` 전수 치환 |
| 006 | FlatRemovalAndGoginActivation | ✅ | Flat mode 제거 + gogin.Generate stub 해소 (MainGenInput struct 포함) |
| 007 | ReactHurlActivation | ✅ | react.Generate + hurl.Generate stub 해소. `adapt_policies.go` (rego→internal policy 어댑터) |
| 008 | OrchestratorWiring | 🟡 | gen_glue.go 배선 교체 완료. genSSaC/genSTML 보류 |
| 009 | StructuralCleanup | 🟡 | generateMain struct 수렴. 나머지 템플릿 분리는 Phase010 과 통합 |
| **010** | **ToulminPoints** | **대기** | **설계 문서 체크포인트 선행 필수 — 다음 세션** |
| 011 | DummyRegressionValidation | 대기 | Tier 1~3 검증 + 구조 건전성 지표 |

## 검증 상태

- `go build ./pkg/... ./internal/... ./cmd/...` ✅
- `go vet` ✅
- `go test ./pkg/...` 전수 통과 ✅
- `fullend validate dummys/gigbridge/specs` 정상 동작 ✅
- `fullend gen dummys/gigbridge/specs /tmp/x` 성공 — **pkg/generate 경유** (glue-gen ✓) ✅
  - 산출: `backend/internal/{api, auth, authz, billing, middleware, model, service, states}`

## 핵심 아키텍처 결정

- **pkg/fullend.Fullstack** = SSOT 파싱 결과 canonical 컨테이너
- **pkg/rule.Ground** = validate + crosscheck + generate 공유 조회 계층
  - 기존 평탄: `Lookup/Types/Pairs/Schemas` (validate 주도)
  - 신규 구조: `Models/Tables/Ops/ReqSchemas` (generate 주도)
- **복사 방식** — internal 은 참조용으로 유지, 별도 Phase 에서 일괄 삭제
- **기술 스택 단일** — `gogin/react` 디렉토리는 수요 기반 확장 슬롯 (nisabit.com 운영 데이터 기반)

## 다음 세션 시작 지점

**Phase010 — Toulmin 포인트 3군데 도입**

설계 선행 필수. 3개 설계 문서 작성 → 사용자 리뷰 → 구현 착수.

대상:
1. `method_from_iface.go` 7-case switch (축 5개)
2. `main` 초기화 블록 조합 (6축 독립)
3. `hurl` 시나리오 순서 (5-phase + topological)

설계 문서 위치:
- `plans/v5/Phase010-MethodDispatchDesign.md`
- `plans/v5/Phase010-MainInitDesign.md`
- `plans/v5/Phase010-ScenarioOrderDesign.md`

## 보류된 후속 작업 (Phase010 이후)

1. **orchestrator 완전 전환** — genSSaC/genSTML 이 여전히 internal 생성기 호출. `parsed.go` 의 SymbolTable 필드도 유지 중. 완전 분리는 별도 Phase.
2. **템플릿 분리** — `main_template.go`, `query_opts_template.go` 의 `text/template` 화. Phase010 main init Toulmin 재설계와 통합 예정.
3. **internal/* 일괄 삭제** — pkg 안정화 후 별도 Phase.
4. **funcspec/policy/statemachine/manifest** 중 internal 에 남은 의존 정리.

## 주요 산출 변경

- `pkg/fullend/` — Fullstack 타입, ParseAll, DetectSSOTs, SSOTKind 등
- `pkg/parser/iface/`, `pkg/parser/sqlc/` 신설
- `pkg/rule/ground_types.go` — 11개 신 타입
- `pkg/ground/populate_models|tables|ops|request_schemas.go` — Ground 채움
- `pkg/contract/` — internal/contract 이식
- `pkg/generate/` — 전체 계층 구조 (gogin/react/hurl + ssac/stml 서브)
- `pkg/rule/model_ref_exists.go` — S-48 용 ModelRefExists warrant

## 참고 문서

- `Phase002-GroundDesign.md` — Phase002 설계 리뷰 산출물 (SymbolTable 접근 패턴 전수 조사 + Ground 설계)
- `internal/gen/*/README.md` 9개 — internal 코드젠 분석 문서 (Phase001 이전에 작성)

## CLAUDE.md 규약 리마인더

- **.ffignore 수정 금지**
- **Co-Authored-By 금지**
- **버전 자동 bump** — `go install -ldflags "-X main.Version=v0.1.N"` (N 증가)
- **커밋+푸시 시 민감 정보 확인 필수**
