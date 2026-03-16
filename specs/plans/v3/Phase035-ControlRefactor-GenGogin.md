# Phase035: dimension 부착 + Q1/Q3 리팩토링 — gen/gogin ✅ 완료

## 목표

`internal/gen/gogin/`의 filefunc 위반 **96건(ERROR 92 + WARNING 4) → 0건**.

## 현황 (filefunc validate 실측)

| 위반 | 건수 | 설명 |
|---|---|---|
| A15 | 56 | `control=iteration`에 `dimension=` 누락 |
| Q1 depth 3 | 15 | 총 nesting depth 초과 |
| Q1 depth 4 | 15 | 〃 |
| Q1 depth 5 | 1 | `collectModelIncludes` |
| Q1 depth 6 | 1 | `generateServerStruct` |
| A9 | 2 | `control=` 누락 (`generateMethodFromIface`, `generateStateMachineSource`) |
| A12 | 2 | sequence인데 iteration 존재 (동일 2파일) |
| Q3 WARNING | 4 | 줄 수 초과 |
| hint | 4 | backtick 템플릿 분리 권장 (범위 제외) |

## 설계

### 1단계: `dimension=` 어노테이션 부착 (56파일)

dimension = 함수 내 for-range 체인의 최대 중첩 수.

**dimension=1 (30파일, Q1 상한 2):**

`attach_directives_in_dir`, `attach_service_directives`, `attach_tsx_directives`, `authz`, `build_file_to_operation_id`, `collect_subscribers`, `convert_path_params_gin`, `extract_base_where`, `finish_query`, `generate_auth_stub_with_domains`, `generate_domain_handler`, `generate_issue_token`, `generate_middleware`, `generate_model_file`, `generate_model_impls`, `generate_refresh_token`, `generate_scan_func`, `generate_server_struct_with_domains`, `generate_state_machines`, `generate_verify_token`, `get_str_slice`, `has_bearer_scheme`, `has_domains`, `hash_claim_defs`, `inject_file_directive`, `parse_return_types`, `sorted_claim_fields`, `transform_service_files`, `transform_source`, `unique_domains`

**dimension=2 (22파일, Q1 상한 3):**

| 함수 | 순회 대상 |
|---|---|
| `build_ownerships_literal` | policies→ownerships |
| `collect_cursor_specs` | paths→operations |
| `collect_funcs` | funcs→sequences |
| `collect_funcs_for_domain` | funcs→sequences |
| `collect_models` | funcs→sequences |
| `collect_models_for_domain` | funcs→sequences |
| `collect_seq_types` | funcs→sequences |
| `domain_needs_db` | funcs→sequences |
| `domain_needs_jwt_secret` | funcs→sequences |
| `generate_central_server` | paths→operations |
| `generate_main` | paths→operations |
| `generate_main_with_domains` | paths→operations |
| `generate_types_file` | paths→operations |
| `has_auth_sequence` | funcs→sequences |
| `has_publish_sequence` | funcs→sequences |
| `infer_bool_states` | states→states |
| `infer_field_type` | states→states |
| `parse_ddl_files` | entries→lines |
| `parse_models_gen` | lines→matches |
| `parse_query_files` | entries→lines |
| `resolve_includes` | specs→columns |
| `resolve_success_status` | paths→operations |

**dimension=3 (3파일, Q1 상한 4):**

| 함수 | 순회 대상 |
|---|---|
| `domain_needs_auth` | funcs→sequences→args |
| `generate_server_struct` | paths→operations→pathParams |
| `transform_service_files_with_domains` | domains→entries→files |

**dimension=4 (1파일, Q1 상한 5):**

| 함수 | 순회 대상 |
|---|---|
| `collect_model_includes` | paths→operations→specs→existing |

### dimension 부착만으로 해소되는 Q1 (12건)

dimension 부착 후 depth ≤ dimension+1이면 Q1 자동 해소:

| 함수 | depth | dimension | Q1 상한 | 결과 |
|---|---|---|---|---|
| `collect_seq_types` | 3 | 2 | 3 | ✅ 해소 |
| `domain_needs_db` | 3 | 2 | 3 | ✅ 해소 |
| `domain_needs_jwt_secret` | 3 | 2 | 3 | ✅ 해소 |
| `generate_types_file` | 3 | 2 | 3 | ✅ 해소 |
| `has_auth_sequence` | 3 | 2 | 3 | ✅ 해소 |
| `has_publish_sequence` | 3 | 2 | 3 | ✅ 해소 |
| `infer_bool_states` | 3 | 2 | 3 | ✅ 해소 |
| `infer_field_type` | 3 | 2 | 3 | ✅ 해소 |
| `parse_ddl_files` | 3 | 2 | 3 | ✅ 해소 |
| `transform_service_files_with_domains` | 3 | 3 | 4 | ✅ 해소 |
| `domain_needs_auth` | 4 | 3 | 4 | ✅ 해소 |
| `collect_model_includes` | 5 | 4 | 5 | ✅ 해소 |

### 2단계: early-continue 적용 (잔여 Q1)

dimension 부착 후에도 Q1 위반이 남는 파일. `if cond { ... }` → `if !cond { continue }` 변환.

**dim=1, depth 3 → depth ≤ 2 (5건):**

| 함수 | depth | 패턴 |
|---|---|---|
| `attach_service_directives` | 3 | for→if IsDir→else if 체인 |
| `finish_query` | 3 | if 조건 반전으로 for 감싸기 해제 |
| `generate_model_file` | 3 | for→if 조건 병합 |
| `inject_file_directive` | 3 | for→if→strings 조작 |
| `generate.go` (sequence) | 3 | sequence 내 블록 추출 필요 |

**dim=2, depth 4 → depth ≤ 3 (12건):**

| 함수 | depth | 패턴 |
|---|---|---|
| `collect_cursor_specs` | 4 | for→for→if→if 병합 |
| `collect_funcs` | 4 | for→for→if→if 병합 |
| `collect_funcs_for_domain` | 4 | for→if domain→for→if 체인 |
| `collect_models` | 4 | for→for→if→if 병합 |
| `collect_models_for_domain` | 4 | for→if domain→for→if 체인 |
| `generate_central_server` | 4 | for→for→if early-continue |
| `generate_main` | 4 | for→for→if 병합 |
| `generate_main_with_domains` | 4 | for→for→if 병합 |
| `parse_models_gen` | 4 | for→for→if 체인 |
| `parse_query_files` | 4 | for→for→if 체인 |
| `resolve_includes` | 4 | for→for→if 체인 |
| `resolve_success_status` | 4 | for→for→if 체인 |

**추출 필요 (early-continue만으로 부족):**

| 함수 | depth | dim | Q1 상한 | 수정 |
|---|---|---|---|---|
| `generate_server_struct` | 6 | 3 | 4 | `collectPathParams()` + 내부 블록 추출 |
| `transform_source` | 4 | 1 | 2 | for 내부 블록을 `rewriteImports()` 등 추출 |

### 3단계: A9/A12 mixed control 분리 (2파일)

| 함수 | 줄 | 현재 | 위반 | 수정 |
|---|---|---|---|---|
| `generateMethodFromIface` | 298 | `control=` 없음 | A9+A12+Q1(4)+Q3(285) | switch 이전 param reorder for를 `reorderCallArgs()` 추출 + switch 각 case 내 for(includes)를 `writeIncludeLoads()` 추출. 본체 switch만 → `control=selection`, Q3 상한 300 이내 |
| `generateStateMachineSource` | 76 | `control=` 없음 | A9+A12 | transition for를 `writeTransitions()` 추출. 본체 → `control=sequence` |

### 4단계: Q3 초과 분리 (4건 WARNING)

| 함수 | 줄 | control | Q3 상한 | 수정 |
|---|---|---|---|---|
| `generateCentralServer` | 137 | iteration | 100 | `writeRoutes()` 추출 (~43줄). 본체 ~94줄 |
| `generateServerStruct` | 137 | iteration | 100 | 2단계 추출과 동시 해결 |
| `generateModelFile` | 104 | iteration | 100 | method 생성 루프(92-115)를 `writeModelMethods()` 추출 |
| `generateMethodFromIface` | 285 | selection(변경 후) | 300 | 3단계에서 해소 |

## 변경 파일

- `dimension=` 추가: 56개 iteration 파일
- `control=` 부착: 2개 (A9 해소)
- early-continue 적용: ~17개 파일
- 함수 추출 신규 파일: ~8개 (`reorder_call_args.go`, `write_include_loads.go`, `write_transitions.go`, `collect_path_params.go`, `rewrite_imports.go`, `write_routes.go`, `write_model_methods.go` 등)

## 검증

1. `go test ./internal/gen/gogin/...`
2. `filefunc validate` — gen/gogin ERROR=0, WARNING=0
3. `fullend gen specs/dummys/zenflow-try05/` — 코드젠 결과 비교
4. `fullend gen specs/dummys/gigbridge-try02/` — 코드젠 결과 비교
