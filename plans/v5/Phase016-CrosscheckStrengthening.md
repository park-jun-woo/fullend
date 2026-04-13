# ✅ Phase016 — CrosscheckStrengthening (완료: e870181, v0.1.10)

> Phase017 실측에서 **crosscheck 가 잡지 못한** 정합성 위반 6종을 `pkg/crosscheck/` 신규 규칙으로 추가. 재발 방지.

## 배경

Phase017 까지 노출된 결함 중 crosscheck 규칙 공백으로 조기 검출 실패한 것들:

| 결함 유형 | 실제 증상 | 현 crosscheck 결과 |
|---------|---------|---------------|
| DDL CHECK vs INSERT seed 불일치 | `psql -f users.sql` 실행 시 CHECK 위반 | 탐지 못 함 |
| DEFAULT FK vs seed row 부재 | INSERT 시 FK 위반 | 탐지 못 함 |
| claims 타입 ↔ DDL 컬럼 타입 | 생성 코드 컴파일 실패 (string vs int64) | 탐지 못 함 |
| SSaC 하드코딩 role ↔ OPA 정책 | 런타임 403 (member 가 admin 액션 호출) | 탐지 못 함 |
| `@empty` 대상 반환 타입 nilable | 컴파일 실패 (`value == nil`) | 탐지 못 함 |
| `@call` 인자 타입 funcspec 호환 | 컴파일 실패 ([]Action vs []ActionInput) | 탐지 못 함 |

## 목표

위 6종을 `fullend validate` 가 **spec 단계에서** 감지. 런타임/컴파일로 가기 전에 ERROR/WARN.

## 전제

- Phase013 완료 ✅
- Phase017/018 와 **독립**. 본 Phase 를 **먼저** 하면 기존 결함 (gigbridge users.sql, DEFAULT FK 등) 이 validate 로 감지돼 Phase017 의 조사 리스트로 직접 활용 가능
- Phase018 (auto seed) 나중에 완료되면 2-2 규칙의 false positive 튜닝 필요

---

## 범위

### 2-1. X-NN: DDL CHECK vs INSERT seed 일치 — **ERROR**

**조건**: 같은 DDL 파일 내 `CREATE TABLE ... CHECK (col IN ('a','b','c'))` 와 `INSERT INTO ... (col) VALUES ('X')` 가 공존하면 `X ∈ {'a','b','c'}` 검증.

**실제 위반**: gigbridge `users.sql` 의 `INSERT ... role='system'` vs `CHECK (role IN ('client','freelancer','admin'))`.

**구현 위치**: `pkg/crosscheck/check_ddl_check_vs_insert.go`

**Severity**: ERROR (의미적 모순, psql 적용 즉시 실패).

### 2-2. X-NN: DEFAULT N FK 컬럼과 seed row 존재 — **WARN**

**조건**: 컬럼이 `DEFAULT <N> REFERENCES <table>(id)` 이고 N 이 정수 리터럴이면, 같은 DDL 에 `<table>` 의 `INSERT ... VALUES (N, ...)` 또는 Phase018 의 auto seed 가 주입되는지 검증.

**실제 위반**: gigbridge `gigs.freelancer_id BIGINT NOT NULL DEFAULT 0 REFERENCES users(id)` + nobody seed 실패.

**구현 위치**: `pkg/crosscheck/check_default_fk_seed.go`

**Severity**: WARN (런타임에만 터짐, 생성 산출물 빌드는 통과).
**Phase018 후**: auto seed 가 on 이면 자동 충족, off 이면 WARN 유지.

### 2-3. X-NN: claims 타입 ↔ DDL 컬럼 타입 일치 — **ERROR**

**조건**: `fullend.yaml` 의 `auth.claims.X: colName[:type]` 가 DDL 의 `users.colName` 컬럼을 매핑한다고 가정하면, 명시된 `type` (또는 기본 string) 과 컬럼 Go 타입이 호환 검증.

**실제 위반**: zenflow 초기 `ID: user_id` (기본 string) vs `users.id BIGINT` (int64).

**구현 위치**: `pkg/crosscheck/check_claims_vs_ddl.go`

**Severity**: ERROR (생성 코드 컴파일 실패).

### 2-4. X-NN: SSaC 하드코딩 role ↔ OPA 정책 role 집합 — **WARN**

**조건**: SSaC 파일에서 `Role: "X"` 하드코딩 (특히 register 계열) 이 발견되면:
- X 가 `fullend.yaml` 의 `roles:` 목록에 포함
- X 가 OPA Rego 의 `input.claims.role == "X"` 분기 중 하나 이상에 등장 (또는 role 무관 allow)

**실제 위반**: zenflow `register.ssac` 의 `Role: "member"` vs OPA 가 `admin` 만 허용.

**구현 위치**: `pkg/crosscheck/check_ssac_role_vs_policy.go`

**Severity**: WARN (런타임에만 권한 거부).

### 2-5. X-NN: `@empty` 대상의 반환 타입 nilable — **ERROR**

**조건**: SSaC `@empty <var>` 의 `var` 가 `@call <Type> <var> = <fn>(...)` 또는 `@get <Type> <var> = ...` 로 바인딩된 경우:
- `<fn>` (funcspec) 의 반환 타입이 pointer/slice/map/interface/chan 중 하나
- 또는 `<Type>` 이 `*Type` 형식

**실제 위반**: zenflow `@empty credits` 인데 `billing.CheckCredits` 반환 타입이 `CheckCreditsResponse` (value).

**구현 위치**: `pkg/crosscheck/check_empty_nilable.go`

**Severity**: ERROR (생성 코드 `if credits == nil` 컴파일 실패).

### 2-6. X-NN: SSaC `@call` 인자 타입 funcspec 호환 — **ERROR**

**조건**: `@call fn({Field: expr})` 에서 각 `Field`/`expr` 쌍의 Go 타입이 funcspec Request struct 의 Field 타입과 호환.

- 단순 값: 타입 equal
- slice: 요소 타입 equal (named 타입 포함)
- struct: 필드 이름·타입 equal (덕 타이핑 금지)

**실제 위반**: zenflow `worker.ProcessActions({Actions: actions})` 에서 `actions [] model.Action` vs funcspec `[]ActionInput`.

**구현 위치**: `pkg/crosscheck/check_call_arg_types.go`

**Severity**: ERROR (컴파일 실패).

---

## 작업 순서

### Step 1. 기존 rule X-번호 체계 확인

`pkg/rule/` 또는 `pkg/crosscheck/` 의 기존 rule 목록 정리. 다음 번호 = N. 신규 6개는 `X-N+1 ~ X-N+6`.

### Step 2. 규칙별 구현 (2-1 → 2-6 순)

각 규칙마다:
- warrant 함수 (`pkg/rule/ruleX.go`)
- crosscheck 검사 함수 (`pkg/crosscheck/check_X.go`)
- 단위 테스트 (`pkg/crosscheck/test_X_test.go` — 양성/음성 샘플)

Phase010 기준점 (2-depth): 단일 axis 검사는 if-else, 다축 조건은 Toulmin warrant.

대부분 규칙은 **2-depth 이내** → 순수 if-else 로 충분. 2-6 (타입 재귀 호환) 만 Toulmin 고려.

### Step 3. 실증 결함으로 양성 테스트

```bash
# Phase017 수정 전 상태를 git stash 나 별도 브랜치로 복원
# 각 규칙이 실제 결함을 감지하는지 확인
```

또는 테스트 안에 "결함 재현 spec" 고정 샘플 포함.

### Step 4. dummy 회귀 확인

Phase017/018 수정된 gigbridge/zenflow 는 새 규칙 하에서도 `fullend validate` ERROR 0.

### Step 5. 리포트 + 커밋

- `reports/crosscheck-rules-phase018.md` — 신규 6규칙 설명 + 양성/음성 예
- 커밋: `feat(crosscheck): 정합성 규칙 6종 추가 (X-N+1 ~ X-N+6)`

---

## 주의사항

### R1. Severity 설정 보수적으로

ERROR 는 정책 오남용 가능. 애매하면 WARN.

- **ERROR**: spec 상 논리 모순이거나 생성 코드 확정 실패
- **WARN**: 런타임에만 실패하거나 사용자 의도일 가능성

### R2. False Positive 검사

각 규칙의 "정당한 예외" 케이스:
- 2-1: INSERT 가 runtime 의도된 테스트 데이터 — 하지만 spec 의 INSERT 는 DB 초기화 의미라 false positive 가능성 낮음
- 2-3: claim 이 DDL 과 무관한 (외부 시스템 발급) 경우 — fullend.yaml 에 `no_ddl_mapping: true` 플래그 지원 고려
- 2-4: role 이 정책 없이 "기본 권한" 만 쓰는 경우 — 검증 건너뛰기
- 2-6: interface 타입 수용 — 구조적 매치 허용

### R3. Phase018 종속 (2-2)

Phase018 auto seed 활성 시 2-2 결함이 런타임에 안 터짐 → 규칙이 "정당한 DEFAULT" 를 WARN 으로 잡는 false positive 될 수 있음. Phase018 완료 후 규칙 튜닝.

### R4. 신규 번호 체계 일관성

기존 X-1 ~ X-72 와 자연스럽게 이어지도록. crosscheck rule manifest (있다면) 업데이트.

---

## 완료 조건 (Definition of Done)

- [x] 2-1 **X-78** DDL CHECK vs INSERT seed — ERROR, 양성 1건 (gigbridge users.sql role='system')
- [x] 2-2 **X-79** DEFAULT FK vs seed row — WARNING, 양성 테스트 (nobody seed 임시 제거 후 1건)
- [x] 2-3 **X-74** claims 타입 vs DDL 컬럼 — ERROR, 양성 2건 (zenflow ID/OrgID 임시 복원 시)
- [x] 2-4 **X-76** SSaC role vs OPA 정책 — WARNING, 양성 1건 (zenflow Register Role=member 임시)
- [x] 2-5 **X-75** `@empty` vs 반환 타입 nilable — ERROR, 양성 2건 (billing value 반환 임시)
- [x] 2-6 **X-77** `@call` 인자 타입 funcspec 호환 — ERROR, 양성 1건 (worker ActionInput 임시)
- [x] dummy 음성 테스트: 정상 상태에서 신규 규칙 X-74~79 중 gigbridge 1건(X-78), zenflow 0건 — **의도된 노출**
- [x] `go build / vet / test ./pkg/...` 통과
- [x] `filefunc validate` 신규 파일 위반 **0** (baseline 37 유지)
- [ ] `reports/crosscheck-rules-phase016.md` 생성 (커밋 시 작성)
- [ ] 커밋: `feat(crosscheck): 정합성 규칙 6종 추가 (X-74~X-79) + ddl/funcspec 파서 확장`

### 부산물

- `pkg/parser/ddl/` — `Table` 에 `Defaults`/`Seeds` 필드 추가. `extract_inserts.go`, `apply_default.go` 등 5파일 신설
- `pkg/parser/funcspec/` — `FuncSpec.ResponsePointer bool` 추가. `first_result_is_pointer.go`, `process_func_decl.go` 신설
- `pkg/crosscheck/` — 6개 rule 파일 + 10개 helper 파일 신설

### 영향

Phase016 시행 후 **gigbridge gen 은 X-78 ERROR 로 차단**. 이는 의도된 동작 — 원래 spec 의 숨어있던 결함이 이제 표면화. **Phase017 에서 gigbridge `users.sql` 의 `role='system'` → 유효 값으로 수정** 시 해소.

## 의존

- Phase013 완료 ✅
- Phase017/018 와 **독립** — 본 Phase 가 "검증 먼저" 트랙의 첫 단계. 기존 결함을 대상으로 규칙 동작 (양성 검증) 을 직접 수행할 수 있다는 이점.

## 다음 Phase

- **Phase017** — RuntimeBugFixing (본 Phase 가 드러낸 결함 + validate 로 못 잡는 순수 런타임 버그 수정)
- **Phase018** — DDLPipelineIntegration (auto seed 활성 시 본 Phase 의 2-2 규칙 의미 조정)
