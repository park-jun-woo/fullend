# Phase009 — 구조 정리 (순수 리팩토링)

## 목표

Phase008 까지 pkg/generate 가 동작하는 상태. 이제 **구조 건전성을 수치로 끌어올리는 순수 리팩토링**.

1. **매개변수 비대 해소** — `generateMain` 등 다매개변수 함수를 struct 로.
2. **결정 분산 해소** — Queue init 등 판정 로직을 한 곳에 취합.
3. **거대 템플릿 처리** — `main_template.go` 같은 큰 문자열 리터럴을 `text/template` 로 분리.

**Toulmin 포인트 도입은 Phase010** (별도 설계 체크포인트 필요).

성공 기준:
- 행동 보존 — 생성 산출물 기능 유지 (bit-동일 의무 없음; 구조적 개선 허용)
- `go build` + `go test ./pkg/...` 통과
- 지표 개선 (Phase011 에서 측정)

---

## 전제

- **Phase008 부분 완료** — pkg/generate 가 실제 동작 (gen_glue 경유).
- internal/gen 은 dead code 로 남음 (본 Phase 에서도 유지).

---

## 원칙

- **순수 리팩토링**: 로직 재구성은 OK, 실제 동작 바뀌면 안 됨
- 한 커밋에 한 가지 개선만
- 회귀 방지를 위해 매 커밋 빌드·테스트 통과

---

## 범위

### 포함

1. **매개변수 비대 해소**
   - `generateMain` (8 params) → `generateMain(input MainGenInput)`
   - 기타 7+ params 함수 전수 조사 후 변환

2. **결정 분산 해소**
   - Queue init 판정 (현재 `generate_main.go` + `collect_subscribers.go` + `has_publish_sequence.go` + `build_queue_blocks.go` 로 분산됨) → `queue_decision.go` 한 곳에 취합
   - 결과 struct 반환: `QueueNeed { Import string; InitBlock string; SubscribeBlock string }`
   - 호출자는 struct 필드만 참조

3. **거대 템플릿 분리**
   - `main_template.go`, `query_opts_template.go` 등 큰 문자열 리터럴
   - Go `text/template` 기반 `.tmpl` 파일 + 렌더 함수로 분리 (단순 템플릿에 한정)
   - 복잡한 분기가 있는 것은 Phase010 Toulmin 대상으로 보류

### 포함하지 않음

- **Toulmin 포인트 3군데 도입** — Phase010 (별도 설계 필요)
- Dummy 검증 — Phase011

---

## 작업 순서

### Step 1. 대상 식별

매개변수 8+ 함수, 결정 분산 패턴, 거대 템플릿을 전수 조사. `grep + wc` 활용 또는 스크립트.

### Step 2. 매개변수 비대 해소 (커밋 단위별)

각 함수 별 커밋:
- `refactor(generate): generateMain 을 MainGenInput struct 로 수렴`
- (기타 발견된 함수 동일 패턴)

### Step 3. 결정 분산 해소

- `queue_decision.go` 신설
- 분산 판정 통합
- 호출자 수정
- 기존 분산 함수 중 사용처 0 된 것 제거

### Step 4. 템플릿 분리 (단순한 것만)

- 큰 문자열 리터럴을 `.tmpl` 별도 파일로 분리
- 렌더 함수 추가
- 빌드·실행 검증

### Step 5. 검증

- `go build ./pkg/... ./internal/... ./cmd/...` 통과
- `go vet` 통과
- `go test ./pkg/...` 통과
- `fullend gen dummys/gigbridge/specs /tmp/x` 실행 에러 없음

---

## 주의사항

### R1. 매 커밋 빌드 통과

리팩토링은 항상 점진. 한 커밋에 여러 개선을 섞지 말 것.

### R2. 행동 변화 최소

- 내부 구조 자유롭게 개선
- 생성된 산출물의 **실질 기능 동일** 유지
- 산출물의 공백·순서 등 사소한 차이는 허용

### R3. "복잡한 것" 은 Phase010 으로

템플릿 분리 중 분기가 많아지면 Phase010 Toulmin 대상. 여기선 단순한 것만.

---

## 완료 조건 (Definition of Done)

- [ ] 매개변수 8+ 함수가 pkg/generate 에 존재하지 않음
- [ ] Queue init 판정이 단일 파일에 집약 (`queue_decision.go` 또는 유사)
- [ ] 거대 템플릿 리터럴 1개 이상 `text/template` 분리
- [ ] `go build` + `go vet` + `go test ./pkg/...` 통과
- [ ] `fullend gen dummys/gigbridge/specs /tmp/x` 실행 성공
- [ ] 커밋 메시지: `refactor(generate): 구조 정리 (매개변수/결정분산/템플릿)`

---

## 다음 Phase

- **Phase010** — Toulmin 포인트 3군데 도입 (설계 체크포인트 포함).
- **Phase011** — Dummy 실용 검증 + 구조 건전성 지표 측정.
