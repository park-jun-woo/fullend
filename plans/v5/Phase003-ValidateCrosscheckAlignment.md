# Phase003 — validate/crosscheck 를 확장된 Ground 에 맞춰 정리

## 목표

Phase002 에서 `pkg/ground/` 에 추가된 **구조적 신 필드** (`Models`, `Tables`, `Ops`, `DTOs` 등) 를 validate/crosscheck 가 활용하도록 점진적으로 마이그하고, **중복된 평탄 populate 를 제거** 한다.

**행동 변화 금지** — 이 Phase 전후로 validate/crosscheck 의 에러 메시지, 검증 결과, 테스트 출력이 bit-동일해야 한다.

---

## 전제

- **Phase001 완료** — pkg/fullend 분리, Feature 리네임.
- **Phase002 완료** — `pkg/ground/` 에 `Models`, `Tables`, `Ops`, `DTOs` 등 구조적 필드와 대응 populate 존재.
- validate/crosscheck 는 여전히 기존 평탄 필드(`Lookup`, `Types`, `Pairs`, `Schemas`)를 소비 중.
- 기존 평탄 populate 와 신 구조적 populate 가 **동일 정보를 두 형태로** 보유 (Phase002 의 용인된 중복).

---

## 범위

### 포함

- validate/crosscheck 의 Ground 소비 지점 전수 조사.
- 신 구조적 필드로 대체 가능한 호출을 점진 치환.
- 사용처가 0 이 된 평탄 populate 는 제거.
- 테스트 업데이트 (필요 시 — 테스트가 Ground 내부 구조에 의존한다면).

### 포함하지 않음

- Ground 에 새 필드 추가 — Phase002 가 담당.
- internal/ 수정 — 본 Phase 범위 밖.
- generator 이식 — Phase004.
- 에러 메시지 개선 등 행동 변화 — 별도 작업.

---

## 작업 순서

### Step 1. 소비 지점 전수 조사

`pkg/validate/**`, `pkg/crosscheck/**` 에서 Ground 접근 패턴 분류:

```bash
grep -rn "g\.Lookup\|g\.Types\|g\.Pairs\|g\.Schemas\|g\.Config\|g\.Vars\|g\.Flags" pkg/validate/ pkg/crosscheck/
```

결과를 다음 범주로 분류:
- **A. 직접 대체 가능** — 신 필드의 구조적 조회로 1:1 대응.
- **B. 간접 대체 가능** — 신 필드에서 파생/필터링 필요.
- **C. 대체 불필요** — Boolean 존재 확인 등 평탄 조회가 본질적으로 더 적합.

범주 A 가 주 대상. 범주 C 는 건드리지 않음.

### Step 2. 치환 단위 결정 및 커밋 분할

소비 지점을 규칙(warrant) 단위로 묶어 **규칙 하나씩 마이그**. 한 번에 한 규칙만 치환 → 그 규칙의 테스트 통과 확인 → 커밋.

예:
- `pkg/validate/ssac/check_model_refs.go` 하나 마이그 → 테스트 → 커밋
- `pkg/crosscheck/check_x9.go` 하나 마이그 → 테스트 → 커밋
- ...

각 커밋:
```
refactor(validate): {규칙명} 을 Ground 구조적 필드로 마이그
```

### Step 3. 사용처 제거된 populate 식별 및 제거

모든 규칙 마이그 후, 기존 평탄 populate 중 사용처 0 이 된 것을 삭제:

```bash
# 예: populate 한 키 검색
grep -rn "g\.Lookup\[\"해당.키\"\]" pkg/validate/ pkg/crosscheck/
# 결과 0 이면 해당 populate 안전 제거 가능
```

제거 후보 (예시 — 실제는 Phase002 설계에 따라 변동):
- `populate_op_params.go` 일부 (→ `populate_ops.go` 로 대체됨)
- 기타 평탄 필드 중 신 구조적 필드로 대체 가능한 것

**비고**: `populate_symbol_table.go` 는 Phase002 Part D 에서 이미 제거 완료됨 (사용처 0건 실측 기반). 본 Phase 에서 별도 처리 불필요.

**실측 필수** — 사용처 확인 없이 제거 금지.

각 제거도 독립 커밋:
```
refactor(ground): populate_xxx.go 제거 — {대체 populate} 로 통합됨
```

### Step 4. 회귀 확인

각 커밋마다:
- `go build ./pkg/...`
- `go vet`
- `go test ./pkg/validate/... ./pkg/crosscheck/... ./pkg/ground/...`

전수 통과 필수.

dummy 프로젝트 기반 종단 확인 (선택, Phase005 이전):
- `go run ./cmd/fullend validate dummys/gigbridge/specs`
- 출력이 Phase002 직후와 bit-동일한지.

### Step 5. 최종 커밋

Phase003 종료 커밋:
```
refactor: validate/crosscheck Ground 구조적 필드 마이그 완료
```

---

## 주의사항

### R1. 행동 보존이 최우선

에러 메시지, 에러 순서, 규칙 적용 순서, 로그 포맷까지 **bit-동일** 해야 한다. 치환 과정에서 "같은 의미지만 다른 출력" 이 발생하면 회귀로 간주.

확인 수단:
- dummy 프로젝트 `fullend validate` 실행 결과를 Phase002 직후 스냅샷과 diff.
- 각 warrant 테스트의 예상 출력이 동일 유지.

### R2. 커밋 세분화 필수

한 번에 여러 규칙을 묶어 마이그하지 말 것. **규칙 하나당 커밋 하나** 원칙:
- 회귀 발생 시 원인 추적이 빠름 (`git bisect` 가능).
- 롤백이 부분적으로 가능 (문제 규칙만 되돌림).
- 코드 리뷰 부담 감소.

### R3. populate 제거는 보수적으로

사용처 0 이 확실한 것만 제거. 확인 방법:
- `grep -rn "\"해당.키\"" pkg/ internal/` — 키 문자열 직접 검색.
- 남아있으면 제거하지 않음.
- generator (Phase004 에서 추가될 예정) 도 고려 — 이 Phase 시점에는 generator 가 pkg 에 없으므로 internal 참조만 확인하면 충분하나, Phase004 착수 전에 재점검 필요.

### R4. 범주 C (대체 불필요) 는 유지

"이 이름 집합에 존재함?" 같은 Boolean 조회는 평탄 `Lookup` 이 구조적 필드보다 간결. 억지로 구조적 필드로 치환하지 않음. Ground 에 **두 인터페이스가 공존** 하는 게 합리적이면 공존 유지.

### R5. 테스트 불변

기존 `pkg/validate/*_test.go`, `pkg/crosscheck/*_test.go` 의 테스트 케이스는 변경하지 않음. 테스트 입력·기대 출력이 그대로인데 내부 구현만 바뀌어서 통과해야 "행동 보존" 의 증거.

### R6. Ground 신 필드 설계 오류 발견 시

Phase002 에서 추가한 Ground 신 필드가 validate/crosscheck 요구를 다 못 담으면 본 Phase 에서 발견됨. 이 경우:
- Phase002 재개 (Ground 필드 재설계) — 소규모라면.
- 또는 현 상태 유지 + 해당 규칙은 기존 평탄 필드 계속 사용 + 별도 이슈 기록.

범위 폭발 방지 위해 후자 권장.

---

## 검증 방법

### 정적 검증

- `go build ./pkg/...` 통과.
- `go vet ./pkg/...` 통과.
- `go test ./pkg/validate/... ./pkg/crosscheck/... ./pkg/ground/...` 전수 통과.

### 동작 검증

- dummy 프로젝트 `fullend validate dummys/gigbridge/specs` 출력이 Phase002 직후와 bit-동일.
- 동일하게 zenflow 에서도 확인.

### 정리 검증

- `grep -rn "g\.Lookup\[" pkg/validate/ pkg/crosscheck/` 결과가 Phase002 대비 감소 (신 필드로 치환된 만큼).
- 사용처 0 populate 파일이 삭제됨.

---

## 완료 조건 (Definition of Done)

- [ ] validate/crosscheck 의 Ground 접근 중 범주 A (직접 대체 가능) 가 전부 신 필드로 마이그됨
- [ ] 사용처 0 이 된 평탄 populate 파일 제거됨
- [ ] 기존 테스트 전부 통과
- [ ] dummy validate 출력 bit-동일
- [ ] 각 마이그 커밋이 규칙 단위로 분리되어 있음
- [ ] Phase003 종료 커밋 메시지 등록

---

## 다음 Phase

- **Phase004** — internal/gen 코드젠을 pkg/generate 로 이식 (복사 방식). generator 가 Ground 신 필드를 직접 소비 (SymbolTable 이식 불필요).
- **Phase005** — dummy 회귀 검증 (baseline 캡처 + diff).
