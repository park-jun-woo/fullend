# Phase036: dimension 부착 + Q1/Q3 리팩토링 — gen/hurl ✅ 완료

## 목표

`internal/gen/hurl/`의 filefunc 위반 **58건(ERROR 57 + WARNING 1) → 0건**.

## 현황 (filefunc validate 실측)

| 위반 | 건수 | 설명 |
|---|---|---|
| A15 | 32 | `control=iteration`에 `dimension=` 누락 |
| Q1 depth 3 | 18 | 총 nesting depth 초과 |
| Q1 depth 4 | 5 | 〃 |
| Q1 depth 5 | 2 | `buildResourceFirstTransition`, `buildTransitionOrder` |
| Q3 WARNING | 1 | `buildScenarioOrder` 130줄 |

## 설계

### 1단계: `dimension=` 어노테이션 부착 (32파일)

**dimension=1 (19파일, Q1 상한 2):**

`can_delete_table`, `collect_required_roles`, `find_matching_capture`, `find_parent_resource`, `find_token_json_path`, `generate_hurl_tests`, `generate_login_body_with_email`, `generate_request_body`, `generate_request_body_with_overrides`, `generate_response_assertions`, `get_response_schema`, `get_str_slice`, `infer_capture_field`, `infer_resource`, `infer_table_from_path`, `match_fk_prefix`, `path_param_util`, `resolve_token_var`, `substitute_path_params`

**dimension=2 (10파일, Q1 상한 3):**

| 함수 | 순회 대상 |
|---|---|
| `build_branch_skip_set` | diagrams→transitions |
| `build_resource_first_transition` | diagrams→transitions |
| `build_scenario_order` | paths→operations |
| `build_state_ops_set` | serviceFuncs→sequences |
| `collect_auth_fk_resources` | authSteps→properties |
| `parse_ddl_files_hurl` | entries→lines |
| `sort_deletes_by_fk` | tables→fkTables |
| `sort_string_slice` | i→j (bubble sort) |
| `topo_sort_delete` | deps→parents |
| `write_auth_section` | paths→operations |

**dimension=3 (3파일, Q1 상한 4):**

| 함수 | 순회 대상 |
|---|---|
| `build_operation_role_map` | policies→rules→actions |
| `build_transition_order` | diagrams→queue→transitions |
| `parse_ddl_check_enums` | entries→matches→values |

### dimension 부착만으로 해소되는 Q1 (9건)

| 함수 | depth | dimension | Q1 상한 | 결과 |
|---|---|---|---|---|
| `build_state_ops_set` | 3 | 2 | 3 | ✅ 해소 |
| `collect_auth_fk_resources` | 3 | 2 | 3 | ✅ 해소 |
| `parse_ddl_files_hurl` | 3 | 2 | 3 | ✅ 해소 |
| `sort_deletes_by_fk` | 3 | 2 | 3 | ✅ 해소 |
| `sort_string_slice` | 3 | 2 | 3 | ✅ 해소 |
| `topo_sort_delete` | 3 | 2 | 3 | ✅ 해소 |
| `write_auth_section` | 3 | 2 | 3 | ✅ 해소 |
| `parse_ddl_check_enums` | 3 | 3 | 4 | ✅ 해소 |
| `build_operation_role_map` | 4 | 3 | 4 | ✅ 해소 |

### 2단계: early-continue 적용 (잔여 Q1)

**dim=1, depth 3 → depth ≤ 2 (8건):**

| 함수 | depth | 패턴 |
|---|---|---|
| `find_matching_capture` | 3 | for→if→if 병합 |
| `find_parent_resource` | 3 | for→if→if 병합 |
| `generate_hurl_tests` | 3 | for→if 조건 반전 |
| `generate_request_body_with_overrides` | 3 | for→if→if 병합 |
| `get_response_schema` | 3 | for→if 조건 반전 |
| `infer_capture_field` | 3 | for→if→if 병합 |
| `path_param_util` | 3 | for→if→if 병합 |
| `substitute_path_params` | 3 | for→if→if 병합 |

**non-A15 파일, depth 3 → depth ≤ 2 (2건):**

| 함수 | depth | 패턴 |
|---|---|---|
| `get_success_http_code` | 3 | 조건 반전 or early-return |
| `write_auth_pair` | 3 | 조건 반전 or early-continue |

**dim=2, depth 4 → depth ≤ 3 (2건):**

| 함수 | depth | 패턴 |
|---|---|---|
| `build_branch_skip_set` | 4 | for→for→if→if 병합 |
| `build_scenario_order` | 4 | for→for→if→if 병합 |

**추출 필요 (early-continue만으로 부족, 4건):**

| 함수 | depth | dim | Q1 상한 | 수정 |
|---|---|---|---|---|
| `find_token_json_path` | 4 | 1 | 2 | 내부 탐색 블록 추출 |
| `generate_response_assertions` | 4 | 1 | 2 | 내부 블록 추출 |
| `build_resource_first_transition` | 5 | 2 | 3 | BFS/갱신 로직 추출 |
| `build_transition_order` | 5 | 3 | 4 | BFS 탐색 로직 추출 |

### 3단계: Q3 초과 분리 (1건 WARNING)

| 함수 | 줄 | control | Q3 상한 | 수정 |
|---|---|---|---|---|
| `buildScenarioOrder` | 130 | iteration | 100 | step 수집/분류 블록 추출 |

## 변경 파일

- `dimension=` 추가: 32개 iteration 파일
- early-continue 적용: ~12개 파일
- 함수 추출 신규 파일: ~5개

## 검증

1. `go build ./internal/gen/hurl/`
2. `filefunc validate` — gen/hurl ERROR=0, WARNING=0
3. `fullend gen specs/dummys/gigbridge-try02/` — Hurl 생성 결과 비교
4. `fullend gen specs/dummys/zenflow-try05/` — Hurl 생성 결과 비교
