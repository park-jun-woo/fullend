# Phase037: dimension 부착 + Q1/Q3 리팩토링 — orchestrator ✅ 완료

## 목표

`internal/orchestrator/`의 filefunc 위반 **65건(ERROR 62 + WARNING 3) → 0건**.

## 현황 (filefunc validate 실측)

| 위반 | 건수 | 설명 |
|---|---|---|
| A15 | 39 | `control=iteration`에 `dimension=` 누락 |
| Q1 depth 3 | 15 | 총 nesting depth 초과 |
| Q1 depth 4 | 4 | `checkPathParamConflicts`, `loadDTOTypes`, `traceArtifacts`, `tracePolicy` |
| Q1 depth 5 | 2 | `findFullendPkgRoot`, `statusCmd` |
| Q1 depth 6 | 1 | `traceFuncSpecs` |
| Q1 depth 7 | 1 | `checkDDLNullableColumns` |
| Q3 WARNING | 3 | `GenWith` 151줄, `TestParseIdempotency` 117줄, `ValidateWith` 102줄 |

## 설계

### 1단계: `dimension=` 어노테이션 부착 (39파일)

**dimension=1 (23파일, Q1 상한 2):**

`chain_cmd`, `count_policy_rules`, `detect_db_engine`, `detect_ssots`, `gen_stml`, `gen_with`, `grep_line`, `inject_func_err_status_from_parsed`, `parsed`, `print_status`, `run_cross_validate`, `sorted_string_keys`, `stml_match_attr`, `to_snake_case`, `trace_ddl`, `trace_ssac`, `trace_states`, `trace_stml`, `validate_funcspec`, `validate_policy`, `validate_states`, `validate_stml`, `validate_with`

> `gen_with`, `run_cross_validate`, `validate_with`는 여러 for-range가 있으나 모두 순차(중첩 아님) → dim=1.

**dimension=2 (15파일, Q1 상한 3):**

| 함수 | 순회 대상 |
|---|---|
| `check_ddl_nullable_columns` | files→lines |
| `check_path_param_conflicts` | paths→segments |
| `check_sqlc_query_duplicates` | entries→scanner |
| `find_ddl_table` | entries→scanner |
| `find_fullend_pkg_root` | dirs→lines |
| `gen_func` | entries→files |
| `load_dto_types` | matches→lines |
| `scan_func_imports` | files→lines |
| `trace_artifacts` | sequences→funcs |
| `trace_func_specs` | callRefs→specs |
| `trace_hurl_scenarios` | paths→operations |
| `trace_openapi` | paths→operations |
| `trace_policy` | policies→rules |
| `validate_ddl` | tables→columns |
| `validate_openapi` | paths→operations |

**dimension=3 (1파일, Q1 상한 4):**

| 함수 | 순회 대상 |
|---|---|
| `status_cmd` | detected→paths→operations |

### dimension 부착만으로 해소되는 Q1 (6건)

| 함수 | depth | dimension | Q1 상한 | 결과 |
|---|---|---|---|---|
| `check_sqlc_query_duplicates` | 3 | 2 | 3 | ✅ 해소 |
| `find_ddl_table` | 3 | 2 | 3 | ✅ 해소 |
| `gen_func` | 3 | 2 | 3 | ✅ 해소 |
| `scan_func_imports` | 3 | 2 | 3 | ✅ 해소 |
| `trace_hurl_scenarios` | 3 | 2 | 3 | ✅ 해소 |
| `trace_openapi` | 3 | 2 | 3 | ✅ 해소 |

### 2단계: early-continue 적용 (잔여 Q1)

**dim=1, depth 3 → depth ≤ 2 (5건):**

| 함수 | depth | 패턴 |
|---|---|---|
| `detect_ssots` | 3 | for→if→if 병합 |
| `gen_with` | 3 | for→if→if 병합 |
| `stml_match_attr` | 3 | for→if→if 병합 |
| `to_snake_case` | 3 | for→if→if 병합 |
| `validate_with` | 3 | for→if/else→if 병합 |

**non-A15, depth 3 → depth ≤ 2 (4건):**

| 함수 | depth | 패턴 |
|---|---|---|
| `determine_module_path` | 3 | 조건 반전 or early-return |
| `run_contract_validate` | 3 | 조건 반전 or early-continue |
| `validate_ssac` | 3 | 조건 반전 or early-continue |
| `parse_idempotent_test` | 3 | 테스트 (F5) — 조건 병합 or 헬퍼 |

**dim=2, depth 4 → depth ≤ 3 (4건):**

| 함수 | depth | 패턴 |
|---|---|---|
| `check_path_param_conflicts` | 4 | for→for→if→if 병합 |
| `load_dto_types` | 4 | for→for→if→if 병합 |
| `trace_artifacts` | 4 | for→for→if→if 병합 |
| `trace_policy` | 4 | for→for→if→if 병합 |

**dim=3, depth 5 → depth ≤ 4 (1건):**

| 함수 | depth | 패턴 |
|---|---|---|
| `status_cmd` | 5 | for→switch→for→for→if 병합. 각 case의 카운트 for를 헬퍼로 추출 가능 |

**추출 필요 (early-continue만으로 부족, 3건):**

| 함수 | depth | dim | Q1 상한 | 수정 |
|---|---|---|---|---|
| `find_fullend_pkg_root` | 5 | 2 | 3 | `isFullendGoMod()` 추출 |
| `trace_func_specs` | 6 | 2 | 3 | `findFuncSpecFile()` 추출 |
| `check_ddl_nullable_columns` | 7 | 2 | 3 | `checkColumnLine()` + `checkSentinelRecord()` 추출 |

### 3단계: Q3 초과 분리 (3건 WARNING)

| 함수 | 줄 | control | Q3 상한 | 수정 |
|---|---|---|---|---|
| `GenWith` | 151 | iteration (dim=1) | 100 | `genAllSSOTs()` 추출 |
| `TestParseIdempotency` | 117 | — (F5 테스트) | 100 | `assertIdempotent()` 헬퍼 추출 |
| `ValidateWith` | 102 | iteration (dim=1) | 100 | 경미 — not-found 블록 추출 or inline 정리 |

## 변경 파일

- `dimension=` 추가: 39개 iteration 파일
- early-continue 적용: ~14개 파일
- 함수 추출 신규 파일: ~8개

## 검증

1. `go test ./internal/orchestrator/...`
2. `filefunc validate` — orchestrator ERROR=0, WARNING=0
3. `fullend validate specs/dummys/zenflow-try05/`
4. `fullend validate specs/dummys/gigbridge-try02/`
