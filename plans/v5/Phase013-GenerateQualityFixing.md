# Phase013 — GenerateQualityFixing

> Phase011 Tier 2 에서 기록한 **zenflow 17건 type-mismatch** 해소. 생성 산출물의 실 빌드·실행 품질 확보.

## 목표

`cd dummys/zenflow/artifacts/backend && go build ./...` 통과. 부수적으로 dummy 빌드 환경 자동화.

기준:
- gigbridge/zenflow 둘 다 `go build ./...` 성공
- 수동 `replace` directive 불필요 (fullend gen 이 자동 주입)
- `go vet` 통과

---

## 범위

### A. type-mismatch 17건 범주화 (Phase011 기록)

Phase011 리포트 (`reports/metrics-phase011.md`) 에서 추출:

#### A1. `currentUser.ID` 타입 불일치 (다수)

- 에러: `cannot use currentUser.ID (variable of type string) as int64 value in struct literal` 및 역방향
- 발생 위치: `internal/service/{auth,action,workflow}/*.go` 다수
- 근본 원인: `currentUser.ID` claim 타입(string) 과 DB primary key(int64) 불일치
- 해결 방향:
  - (a) currentUser claim 에 ID 를 int64 로 선언 (fullend.yaml 에 type hint)
  - (b) 생성기가 타입 변환 코드 자동 삽입 (`strconv.ParseInt`)
  - (c) 생성기가 JWT claim 에서 int64 로 파싱하도록 변경
- 권장: **(a) + (c)** — claim 정의에서 원천 타입 고정, 파싱 시 int64. 생성기는 type 조회 후 대입.

#### A2. `billing.CheckCreditsResponse` pointer/value (2건)

- 에러: `invalid operation: credits == nil (mismatched types billing.CheckCreditsResponse and untyped nil)`
- 발생 위치: `workflow/activate_workflow.go`, `workflow/execute_workflow.go`
- 근본 원인: `CheckCredits` 함수 반환 타입이 value 인데 호출부가 `== nil` 비교
- 해결 방향:
  - (a) 반환 타입을 pointer 로 통일 (`*CheckCreditsResponse`)
  - (b) 호출부가 zero-value 비교하도록 생성
- 권장: **(a)** — 외부 서비스 호출은 pointer 반환이 관례

#### A3. `[]model.Action` → `[]worker.ActionInput` 어댑터 누락 (1건)

- 에러: `cannot use actions (variable of type []model.Action) as []worker.ActionInput value`
- 발생 위치: `workflow/execute_workflow.go`
- 근본 원인: SSaC 가 Action → ActionInput 변환 표현 부재
- 해결 방향:
  - (a) 생성기가 `for _, a := range actions { inputs = append(inputs, worker.ActionInput{...}) }` 자동 삽입
  - (b) SSaC 에 어댑터 표기법 도입 (`@adapt`)
- 권장: **(a) 단기** + **(b) 중기** (어댑터 표기 확장은 별도 Phase)

#### A4. dummy go.mod `replace` 자동 주입

- 현상: 생성된 `dummys/*/artifacts/backend/go.mod` 에 `replace github.com/park-jun-woo/fullend => <local>` 없음
- 해결: `generate_main.go` 또는 관련 go.mod 생성 로직에 `-local-fullend-path <path>` 플래그 추가해 주입. 기본값은 `$FULLEND_LOCAL_PATH` 환경변수.

---

## 작업 순서

### Step 1. 정밀 재현

```bash
rm -rf dummys/zenflow/artifacts
go run ./cmd/fullend gen dummys/zenflow/specs dummys/zenflow/artifacts
cd dummys/zenflow/artifacts/backend
go build ./... 2>&1 > /tmp/zenflow-errors.txt
```

17 에러 전수 기록. 생성 파일 경로·라인 확보.

### Step 2. A1 — currentUser.ID 타입 정합

1. `dummys/zenflow/specs/fullend.yaml` 의 claims 정의에 `id: int64` 명시 (이미 있으면 생성기 조회 로직 점검)
2. `pkg/generate/gogin/` 의 currentUser 참조 코드가 claim 타입을 조회해 dst 타입에 맞춰 cast 삽입
3. 생성 대상:
   - `c.MustGet("currentUser").(auth.CurrentUser).ID` — auth.CurrentUser.ID 가 int64 인지 확인
   - 대입 시 `int64(cu.ID)` 불필요하게 변환 말고 원천 타입 맞추기

### Step 3. A2 — billing response pointer

1. `pkg/generate/gogin/` 의 service 호출 반환 타입 처리 점검
2. response 타입이 외부 call 인 경우 pointer 반환으로 통일
3. 호출부 `if credits == nil { ... }` 가 동작하게

### Step 4. A3 — Action → ActionInput 어댑터

1. SSaC 에 `actions := WorkflowModel.FindActions(ctx, workflowID)` 같은 구문이 있다면 반환형 수집
2. 이후 사용처가 다른 구조체 slice 를 기대하면 생성기가 for-range 어댑터 삽입
3. 초기 구현은 필드 매핑 단순 — field 이름 동일 가정

### Step 5. A4 — go.mod replace 자동 주입

1. `pkg/generate/gogin/generate_main.go` (또는 go.mod 생성 지점) 에서 현재 fullend 로컬 경로 기록
2. 플래그 `--local-fullend-path` 또는 env `FULLEND_LOCAL_PATH` 지원
3. 미지정 시 현재 동작 유지 (사용자가 수동 replace 할 것)

### Step 6. 검증

```bash
go run ./cmd/fullend gen dummys/gigbridge/specs dummys/gigbridge/artifacts
cd dummys/gigbridge/artifacts/backend && go build ./... && cd -

go run ./cmd/fullend gen dummys/zenflow/specs dummys/zenflow/artifacts --local-fullend-path $(pwd)
cd dummys/zenflow/artifacts/backend && go mod tidy && go build ./...
```

둘 다 `exit 0` 이 목표.

### Step 7. (선택) hurl 스모크

```bash
cd dummys/gigbridge/artifacts/backend && go run ./cmd/... &
sleep 2
hurl --test --variable host=http://localhost:8080 ../tests/*.hurl
kill %1
```

Tier 3 수준. 실패해도 통과.

---

## 주의사항

### R1. Phase011 리포트 갱신

fix 후 `scripts/structural_metrics.go` 재실행 대신 Tier 2 결과 섹션만 업데이트 (zenflow 17 → 0).

### R2. A3 의 어댑터는 최소 구현

자동 어댑터 생성은 "필드 이름 동일 가정" 만 지원. 이름 다른 경우나 변환 로직 필요한 경우는 별도 Phase (SSaC `@adapt` 표기 도입) 로 분리.

### R3. A4 replace 주입은 옵트인

기본 동작을 바꾸지 않음. 플래그 제공만.

### R4. 범위 폭주 금지

gigbridge/zenflow 외 dummy 신설이나 스펙 확장은 본 Phase 에서 금지. 기존 17건만 해소.

---

## 완료 조건 (Definition of Done)

- [ ] zenflow 17건 type-mismatch 전부 해소 (`go build ./...` exit 0)
- [ ] gigbridge 빌드 유지 (회귀 없음)
- [ ] `go vet ./pkg/... ./internal/... ./cmd/...` 통과
- [ ] `go test ./pkg/...` 통과
- [ ] go.mod replace 자동 주입 플래그 제공 (옵트인)
- [ ] `reports/metrics-phase011.md` 의 Tier 2 섹션 갱신 또는 `reports/metrics-phase013.md` 신설
- [ ] 커밋: `fix(generate): zenflow type-mismatch 17건 해소 + dummy go.mod replace 옵션`

## 의존

- **Phase012 완료 권장** (pkg 전환 완결 후 생성 품질 수정이 안전)
- 단, 독립 진행 가능 (generate 경로는 이미 pkg)

## 다음 Phase

- **Phase014** — internal/* 일괄 삭제
