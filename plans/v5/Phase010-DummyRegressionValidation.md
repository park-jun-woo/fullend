# Phase010 — Dummy 실용 검증 + 구조 건전성 지표

## 목표

Phase009 결과물이 **실용적으로 동작하고**, **구조 건전성이 internal 대비 개선** 되었는지 지표로 확인한다.

**중요한 재정의**: 이 Phase 는 "internal 과 diff 0" 을 요구하지 않는다. internal 은 참고용이지 재현 대상이 아니다. 성공 기준은:

1. pkg/generate 가 **빌드 수준에서 안정** 하다.
2. 생성 산출물이 **부분적으로라도 동작** 한다.
3. 코드 구조가 **internal 대비 건전** 해졌다.

기능이 일부 불완전해도 구조가 건전하면 OK — 차후 프로덕션화 단계에서 기능을 채운다.

---

## 전제

- **Phase009 완료** — pkg/generate 가 살아있고 orchestrator 가 pkg 경로로 배선됨.
- `internal/gen` 은 복사 방식으로 보존됨 (참조용 — 본 Phase 에서 호출하지 않음).
- `go build ./pkg/...` 과 `go test ./pkg/ground/...` 가 이미 통과된 상태 (Phase005~009 의 DoD).

---

## 검증 수준 (3 tier)

### Tier 1 — 필수 통과 (block if failed)

이게 실패하면 구조가 불안정한 것. 진행 중단.

- `go build ./pkg/... ./internal/... ./cmd/...` 통과.
- `go vet ./pkg/...` 통과.
- `go test ./pkg/ground/...` 통과.

### Tier 2 — 목표 (best effort, block 아님)

- `go build ./...` in `dummys/gigbridge/artifacts/backend/` 통과.
- `go build ./...` in `dummys/zenflow/artifacts/backend/` 통과.
- `hurl --test` 실 서버 기동 후 일부 통과.

### Tier 3 — 허용 (실패해도 OK)

- 핸들러 body 부분 누락 (TODO, 스텁 포함).
- 산출물 파일 일부 누락.
- `hurl --test` 부분 실패.
- 생성 코드 일부에 타입 추론 오류.

Tier 2/3 실패는 리포트에 기록하되 pass 판정 유지. 개선은 차후 Phase 로 이월.

---

## 구조 건전성 지표

Phase004~009 이식·리팩토링 중 "구조 개선" 원칙을 실제 수치로 확인.

### 측정 대상 비교: internal vs pkg

| 지표 | 측정 방법 | 기대 방향 |
|------|---------|----------|
| 평균 함수 매개변수 수 | Go AST 순회, 매개변수 count | pkg < internal |
| 최대 매개변수 함수 개수 | 8+ params 함수 수 | pkg << internal |
| 함수 복잡도 (평균) | cyclomatic complexity | pkg ≤ internal |
| 평균 파일 크기 (줄) | `wc -l` 평균 | 유사 또는 감소 |
| 중복 `*WithDomains` 쌍 | 이름 패턴 탐색 | **pkg = 0** (Flat 제거) |
| Toulmin 적용 지점 | `toulmin.NewGraph` 호출 수 | **pkg ≥ 2** |
| 결정 분산도 (queue init 등) | 특정 판정이 걸친 파일 수 | pkg < internal |

### 측정 도구

간단한 Go 스크립트 `scripts/structural_metrics/` 작성:

```
scripts/structural_metrics/
├── README.md
├── count_params.go           — 함수별 매개변수 수 분포
├── count_toulmin.go          — Toulmin 사용 지점 수
├── count_lines.go            — 파일당 줄 수
├── detect_duplicates.go      — _with_domains 류 중복 패턴
└── compare.go                — internal vs pkg 비교 리포트
```

각 스크립트는 `main()` 하나에 단순 수치 반환. 기존 `scripts/` 의 다른 유틸과 독립 빌드.

### 리포트 형식

```
# Structural Metrics Report (Phase010)

## 매개변수 분포
                   internal   pkg/generate
평균               3.2        2.4   ✓ (24% 감소)
중앙값             2          2
최대                9          4
8+ params 함수     7          0    ✓

## Toulmin 적용
                   internal   pkg/generate
warrant 선언       0          N
graph 선언         0          M   ✓

## 중복 패턴
                   internal   pkg/generate
*WithDomains 쌍    14         0    ✓

## 평균 파일 크기
                   internal   pkg/generate
gogin/ 평균        38.5줄     X줄
ssac generator/    22.0줄     Y줄

## 결정 분산 (Queue init 예시)
                   internal   pkg/generate
관련 파일 수       7          Z    ✓ (if Z < 7)
```

---

## 작업 순서

### Step 1. Tier 1 필수 통과 확인

```bash
go build ./pkg/... ./internal/... ./cmd/...
go vet ./pkg/... ./internal/... ./cmd/...
go test ./pkg/ground/...
```

하나라도 실패하면 Phase004 로 복귀 (해당 문제 수정 후 Phase005 재개).

### Step 2. 산출물 생성

```bash
go run ./cmd/fullend gen dummys/gigbridge/specs dummys/gigbridge/artifacts
go run ./cmd/fullend gen dummys/zenflow/specs  dummys/zenflow/artifacts
```

에러가 나면 기록하고 계속 (Tier 2 허용). 에러 로그를 Phase011+ 의 입력으로 쌓음.

### Step 3. Tier 2 목표 검증

각 dummy artifact 에서:
```bash
cd dummys/gigbridge/artifacts/backend && go build ./... ; cd -
cd dummys/zenflow/artifacts/backend && go build ./... ; cd -
```

실패하면 원인 범주화:
- **import 경로 오류** → 생성기 import 처리 버그
- **타입 불일치** → Ground 소비 로직 결함
- **누락 파일** → 생성기 범위 누락
- **문법 오류** → 템플릿 조립 버그

범주별 카운트만 리포트에 기록. 이 Phase 에서 고치지 않음 (Phase011+ 로 이월).

Hurl 검증 (best effort):
```bash
# 가능하면 실 서버 기동
cd dummys/gigbridge/artifacts/backend && go run ./cmd/... &
sleep 2
hurl --test --variable host=http://localhost:8080 artifacts/tests/*.hurl
# teardown
kill %1
```

실패 케이스 카운트만 기록.

### Step 4. 구조 건전성 측정

```bash
go run ./scripts/structural_metrics/compare.go > reports/metrics-phase005.md
```

리포트 생성. 기대 방향 달성 여부 확인.

### Step 5. 판정

- Tier 1 전부 통과 + 구조 지표 악화 없음 → **Phase010 통과**.
- Tier 1 실패 → Phase006~009 복귀.
- 구조 지표 악화 (예: 매개변수 평균 증가, Toulmin 0건) → Phase004 구조 정리 보강 후 재측정.

### Step 6. 최종 커밋

```
feat(phase005): 실용 검증 통과 — 구조 건전성 internal 대비 개선
```

리포트(`reports/metrics-phase005.md`), 스크립트(`scripts/structural_metrics/*.go`) 포함.

---

## 주의사항

### R1. "불완전 허용" 은 Tier 2/3 에만 적용

Tier 1 (go build, go vet, pkg/ground test) 실패는 **구조 안정성 불안** 의 증거라 반드시 해결. Phase009 완료 기준이 이미 Tier 1 포함.

### R2. 구조 지표 악화 대응

지표 악화(예: 매개변수 평균 증가)가 나오면:
- 어느 함수가 원인인지 추적.
- Phase009 의 즉시 해소 누락분 식별.
- Phase004 로 복귀해 해당 함수 정리 후 Phase010 재측정.

이 Phase 에서 해결하지 않음 (Phase 경계 명확히).

### R3. hurl 검증은 선택

서버 기동이 필요하고 환경 의존적이라 CI 에서 생략 가능. 로컬 개발에서 권장.

### R4. dummy 프로젝트 동결

Phase010 기간 중 gigbridge·zenflow SSOT 를 바꾸지 않음. 변경 시 구조 지표 비교가 오염.

### R5. internal/gen 호출 금지

이 Phase 에선 internal 경로로 gen 하지 않음. pkg/generate 만 대상.
(이원 경로 캡처 기법은 v5 초기 설계에서 폐기됨.)

### R6. 리포트 보존

구조 지표 리포트는 git 에 커밋해 후속 Phase 의 비교 기준선으로 활용.

---

## 의존성

- **Phase009 완료** — 본 Phase 전제.
- **외부 도구**: `go`, `hurl` (선택).
- **측정 스크립트**: 본 Phase 에서 새로 작성 (`scripts/structural_metrics/`).

---

## 완료 조건 (Definition of Done)

- [ ] Tier 1 필수 통과 항목 전부 OK
- [ ] `scripts/structural_metrics/*.go` 작성 + 실행 검증
- [ ] `reports/metrics-phase005.md` 생성 + 커밋
- [ ] 구조 지표가 internal 대비 악화 없음
  - [ ] 평균 매개변수 수: pkg ≤ internal
  - [ ] 8+ params 함수: pkg < internal
  - [ ] Toulmin 적용: pkg ≥ 2곳
  - [ ] `*WithDomains` 중복: pkg = 0
- [ ] Tier 2 결과 기록됨 (통과/실패 관계없이 리포트에 수치)
- [ ] 최종 커밋: `feat(phase005): 실용 검증 통과 — 구조 건전성 internal 대비 개선`

---

## 다음 Phase 예고

- **Phase011+ (프로덕션화)** — Tier 2/3 에서 기록된 실패를 하나씩 해소해 생성 산출물의 실제 동작 수준을 올림. 핸들러 body 채우기, 타입 오류 수정, 누락 기능 추가 등.
- **Phase00N (internal 삭제)** — pkg/generate 가 internal 을 완전히 대체한 시점에 일괄 삭제.
