# ⛔ Phase015 — TemplateAndResiduals (실행 금지 — 현재 진행 안함)

> 문자열 조립 기반 템플릿을 `text/template` 로 전환 + Phase012~014 에서 남은 잔여 internal 의존 정리 + filefunc F1/F2 완화 논의 결론.

## 목표

v5 로드맵의 **미결 장기 항목 3가지** 를 한 Phase 에서 정리:

1. **템플릿 전환** — `main_template.go`, `query_opts_template.go` 등 문자열 조립 → `text/template` 화. 유지보수성 + 테스트 용이성.
2. **잔여 internal 의존** — funcspec / policy / statemachine / manifest 에서 internal 참조가 남아있으면 제거.
3. **filefunc F1/F2 완화 논의** — Phase010 에서 발견된 "응집 atom 강제 분해" 비용을 정식 결정으로.

---

## Part A. 템플릿 전환

### 대상

`pkg/generate/gogin/` 의 `*_template.go` 파일군:

- `main_template.go` — cmd/main.go 전체 조립 (현재 `mainWithDomainsTemplate` 함수가 fmt.Sprintf + 문자열 결합)
- `query_opts_template.go` — QueryOpts 생성
- (조사 필요) 그 외 `*_template.go`

### 작업 방식

- 대상 파일마다 `pkg/generate/gogin/templates/<name>.tmpl` 텍스트 템플릿으로 추출
- `//go:embed templates/*.tmpl` 로 인라인 포함
- 렌더링은 `text/template` API. 데이터는 기존 입력 struct 재사용
- 기존 테스트 + 생성 결과 bit-level 비교 테스트 추가

### 이점

- 템플릿과 조립 로직 분리
- 변경 시 diff 가 템플릿에만 집중
- 템플릿 단독 프리뷰 가능

### 범위 한정

`text/template` 화 대상은 `*_template.go` 접미 파일 군 **만**. `writeXXXMethod` 같은 조립 함수들은 범위 외.

---

## Part B. 잔여 internal 의존 정리

### 조사 대상

Phase014 삭제 후 남은 `internal/` 서브패키지:

- `internal/ssac/` (루트) — parser 등 잔여
- `internal/stml/` (루트) — 동일
- `internal/funcspec/`
- `internal/policy/`
- `internal/statemachine/`
- `internal/manifest/`
- `internal/orchestrator/` (본체)
- `internal/crosscheck/` ← pkg/crosscheck 로 이미 이관됐는지 확인
- `internal/reporter/`
- `internal/projectconfig/`

### 작업

```bash
rg "\"github.com/park-jun-woo/fullend/internal/" pkg/ --glob '*.go'
```

pkg 쪽에서 남은 internal import 전수 제거. 필요하면 pkg 로 이식.

반대 방향:
```bash
rg "\"github.com/park-jun-woo/fullend/pkg/" internal/ --glob '*.go'
```

internal → pkg 의존 허용. 역방향만 금지.

### 목표

- `pkg/` 가 `internal/` 을 import 하지 않음 (순방향 엄격)
- pkg 가 완전한 독립 패키지 집합이 됨

---

## Part C. filefunc F1/F2 완화 논의 결론

Phase010 에서 제기된 비용:
- 응집 atom(Pattern enum + MethodFacts + NewMethodFacts + DecideMethodPattern) 을 **4 파일로 분해** 강제
- Go stdlib 관례(`net/http/request.go` 수십 개 심볼 공존) 와 충돌
- 테스트 파일도 F1 적용 시 Go 테이블 드리븐 관례와 충돌

### 결정 필요 항목

1. **F1/F2 유지** — 엄격 준수. 분해 비용 수용.
2. **`//ff:group` 도입** — 응집 단위를 한 파일에 허용하는 메타. 여러 심볼을 그룹으로 묶어 인덱서가 한 단위로 취급.
3. **테스트 파일 면제** — `_test.go` 는 F1/F2 검사 제외. 기존 baseline 이 이미 이렇게 동작 중.
4. **primary symbol 규칙** — 파일당 1 primary symbol + 관련 helper/type 허용.

### 권장 프로세스

이 Phase 에서 **선택안 결정만** 수행. 구현은 별도 filefunc 리포에서:
- 결정 사항 문서화: `docs/filefunc-policy.md` (본 repo)
- filefunc 리포에 이슈 제출 (필요 시)

본 Phase 범위는 **결정 + 문서화**. 구현은 외부 도구 영역.

---

## 작업 순서

### Step 1. 조사

- Part A: `*_template.go` 파일 목록 + 크기 파악
- Part B: `rg` 로 internal 의존 전수 조사
- Part C: Phase010 경험 회고 + 3가지 선택안 비교 작성

### Step 2. Part A 템플릿 전환

파일별로 1 커밋. template 추출 → embed → 렌더 → 기존 테스트 통과 확인.

### Step 3. Part B 의존 정리

pkg → internal import 0건 달성. 필요 시 pkg 로 이식.

### Step 4. Part C 결정

선택안 평가 + 결정 문서화. 사용자 승인 받아 최종안 반영.

### Step 5. 최종 검증

```bash
go build ./...
go vet ./...
go test ./pkg/...
go run ./cmd/fullend gen dummys/gigbridge/specs /tmp/g && \
  (cd /tmp/g/backend && go build ./...)
go run scripts/structural_metrics.go > reports/metrics-phase015.md
```

---

## 주의사항

### R1. Part 간 독립성

A/B/C 는 독립 진행 가능. 순서 자유. 다만 한 Phase 에서 묶는 이유는 **v5 잔여 정리 통합 마감** 의미.

### R2. Part A 는 행동 보존

템플릿 전환 후 생성 산출물이 **공백·순서 포함 bit-level 동일** 목표. 차이 발생 시 template 정정.

### R3. Part B 가 pkg → internal 순방향 남기면 실패

pkg 가 internal 을 참조하는 한 "pkg 독립" 달성 불가. 0건 필수.

### R4. Part C 결정은 재논의 가능

filefunc 정책 변경 후 baseline 위반 수가 달라지면 Phase 재측정 필요.

---

## 완료 조건 (Definition of Done)

- [ ] Part A: 모든 `*_template.go` 가 `text/template` + embed 로 전환
- [ ] Part A: 생성 산출물 bit-level 동일성 테스트 통과
- [ ] Part B: `rg "\"github.com/park-jun-woo/fullend/internal/" pkg/` 결과 0 hit
- [ ] Part C: `docs/filefunc-policy.md` 작성 + 결정안 확정
- [ ] `go build ./... && go vet ./... && go test ./pkg/...` 통과
- [ ] `fullend gen` + dummy backend 빌드 통과 유지
- [ ] `reports/metrics-phase015.md` 생성
- [ ] 커밋: `feat(generate): template 전환 + pkg 의존 독립 + filefunc 정책 확정`

## 의존

- **Phase014 완료** — internal/gen 등 삭제 후 잔여 internal 이 명확해짐

## 다음 Phase

v5 마감. 이후는 프로덕션 운영 단계 (별도 로드맵).

---

## v5 로드맵 종료 선언

Phase015 통과 시 v5 의 3가지 목표 완결:

1. **구조 건전성** — pkg 가 독립 패키지. internal 대부분 제거 or 역방향 의존만 허용.
2. **결정 로직 수렴** — Decide\* 순수 함수 패턴 + 2-depth 기준점.
3. **실용 동작** — gigbridge/zenflow 가 실제 빌드·구동 가능.

이후 로드맵(v6+) 은 기능 확장(queue 백엔드 다양화, 스토리지 어댑터 확장 등) 이 중심.
