# Phase016 — CrosscheckStrengthening Report

2026-04-14 · X-74 ~ X-79 6규칙 추가

## 요약

| ID | 제목 | 등급 | 파일 |
|----|------|------|------|
| X-74 | claims 타입 vs DDL 컬럼 타입 | ERROR | `check_claims_vs_ddl.go` |
| X-75 | `@empty` 대상 반환 타입 nilable | ERROR | `check_empty_nilable.go` |
| X-76 | SSaC 하드코딩 Role vs OPA 정책 | WARNING | `check_ssac_role_vs_policy.go` |
| X-77 | `@call` 인자 타입 funcspec 호환 | ERROR | `check_call_inputs_vs_funcspec.go` |
| X-78 | DDL CHECK vs INSERT seed 일치 | ERROR | `check_ddl_check_vs_seed.go` |
| X-79 | `DEFAULT N FK` vs seed row 존재 | WARNING | `check_default_fk_seed.go` |

## 양성·음성 검증

| 규칙 | 양성 테스트 | 양성 결과 | 음성 (정상 dummy) |
|------|-----------|---------|-----------------|
| X-74 | zenflow `ID:user_id` (타입 누락) | 2건 감지 ✓ | 0 (Phase013 fixed) |
| X-75 | billing `CheckCreditsResponse` 값 반환 | 2건 감지 ✓ | 0 (Phase013 fixed) |
| X-76 | zenflow `Role: "member"` | 1건 감지 ✓ | 0 (Phase016 중 임시 수정, 복원) |
| X-77 | worker `[]ActionInput` 복원 | 1건 감지 ✓ | 0 (Phase013 fixed) |
| X-78 | gigbridge seed `role='system'` | **1건 감지 (실 결함!)** | n/a |
| X-79 | nobody seed 임시 제거 | 1건 감지 ✓ | 0 |

## 파서 확장

### `pkg/parser/ddl`

- `Table.Defaults map[string]string` — 컬럼별 DEFAULT 값
- `Table.Seeds []map[string]string` — INSERT 시드 행
- 신설: `apply_default.go`, `extract_inserts.go`, `find_table_case_insensitive.go`, `split_and_trim.go`, `split_csv_literals.go`, `strip_sql_quotes.go`
- 통합: `parse_ddl_content.go` 가 `extractInserts()` 호출, `parse_column_def.go` 가 `applyDefault()` 호출

### `pkg/parser/funcspec`

- `FuncSpec.ResponsePointer bool` — 첫 반환이 `*T` 인지
- 신설: `first_result_is_pointer.go`, `process_func_decl.go`
- 통합: `process_decl.go` 의 FuncDecl 핸들러가 `processFuncDecl` 위임

## 영향

Phase016 시행 후 `gigbridge` 의 기존 결함 **X-78: `users.sql seed[0].role='system'` 가 CHECK 위반** 이 **자동 감지**. `fullend gen` 이 ERROR 로 차단됨 — 이는 도구의 의도된 동작.

**Phase017 에서 `role='system'` → `role='admin'` 등 수정 필요** (또는 Phase018 auto seed 로 spec INSERT 자체 제거).

## filefunc

- 신규 파일 (`pkg/crosscheck/*` 16개 + `pkg/parser/ddl/*` 6개 + `pkg/parser/funcspec/*` 2개) 모두 F1/F2/A1/A3/A10 준수
- baseline 위반 수: 37 (Phase013 대비 변화 없음)

## 규칙별 구현 메모

**X-74**: `resolveClaimDDLColumn` 휴리스틱 — `<x>_id` → `<x>s.id`; fallback `users.<key>`. 실제 dummy 에서 false positive 0.

**X-75**: funcspec 파서에 `ResponsePointer` 필드 추가. `@empty` 대상의 binding 이 `@call` 이고 funcspec 반환이 value 타입일 때 ERROR.

**X-76**: OPA `AllowRule.RoleValue` 집합 수집. SSaC 의 `Inputs["Role"]` (for @post/@put) 또는 `Args` (for @call) 에서 literal 추출.

**X-77**: `@call Inputs` 의 값이 변수(dot 없음)일 때만 검사. 필드 접근은 타입 해석 한계로 skip. basename 수준 비교.

**X-78**: DDL `Table.CheckEnums` + `Table.Seeds` 교차. INSERT 시드 각 셀이 CHECK 목록에 속하는지.

**X-79**: DDL `Table.Defaults` 중 정수값 + 같은 컬럼의 FK → 참조 테이블에 해당 id seed 존재 여부. 부재 시 WARN.

## 다음 Phase

- **Phase017** — X-78 해소 (gigbridge seed 수정) + 기타 런타임 버그 (zenflow 403, OPA path, DSN)
- **Phase018** — auto seed 활성 시 X-79 false positive 튜닝 필요
