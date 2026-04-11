# pkg/validate

단일 SSOT 자체 정합성 검증. `pkg/rule` 공통 함수 + 고유 함수로 구성.

## 패키지 구조

| 패키지 | parser 대응 | 규칙 ID | 설명 |
|--------|------------|---------|------|
| `symbol/` | — | — | 이름→타입 조회 허브. DDL/OpenAPI/Model/SQLc에서 구축. 모든 validator가 공유 |
| `manifest/` | `parser/manifest/` | C-1 | fullend.yaml 로드 검증 |
| `ddl/` | `parser/ddl/` | D-1~D-3 | sqlc 중복, NOT NULL, 센티널 레코드 |
| `openapi/` | kin-openapi 직접 | O-1 | path 파라미터 충돌 |
| `ssac/` | `parser/ssac/` | S-1~S-58 | 필수 필드, 변수 흐름, 모델 검증, @subscribe 제약 등 |
| `stml/` | `parser/stml/` | TM-1~TM-12 | fetch/action 바인딩, 파라미터, 컴포넌트 참조 |
| `statemachine/` | `parser/statemachine/` | ST-1 | 파싱 검증 |
| `rego/` | `parser/rego/` | P-1 | 파싱 검증 |
| `funcspec/` | `parser/funcspec/` | F-1 | built-in 패키지명 충돌 |
| `hurl/` | `parser/hurl/` | H-1 | .feature 파일 deprecated |
| `model/` | — | M-1 | model/ 디렉토리 비어있음 |
| `contract/` | — | CT-1~CT-2 | //fullend 디렉티브 검증 |

## Toulmin 매핑

```
claim   = 검증 대상 (ServiceFunc, PageSpec, sequence, field name 등)
ground  = *rule.Ground (심볼 테이블에서 구축)
backing = 규칙별 Backing struct (pkg/rule 공통 또는 고유)

warrant  = 기본 규칙 ("변수가 선언 후 사용되어야 한다")
defeater = 예외 ("currentUser는 암묵적 선언" → IsImplicitVar)
```

## 규칙 → pkg/rule 매핑

### 공통 함수 사용 (pkg/rule)

| 규칙 ID | pkg/rule 함수 | backing 요약 | 패키지 |
|---------|---------------|-------------|--------|
| S-1~S-10, S-12~S-21, S-23~S-24 | `FieldRequired` | SeqType별 Field+Required/Forbidden | `ssac/` |
| S-27~S-30 | `VarDeclared` | Ground.Vars 조회 | `ssac/` |
| S-26 | `NameFormat` | Pattern=`dot-method` | `ssac/` |
| S-31 | `ForbiddenRef` | LookupKey=`ssac.configPrefix` | `ssac/` |
| S-32 | `ForbiddenRef` | LookupKey=`publish.forbidden` | `ssac/` |
| S-33 | `ForbiddenRef` | LookupKey=`ssac.reservedSource` | `ssac/` |
| S-34~S-35 | `ForbiddenRef` | LookupKey=`go.reserved` | `ssac/` |
| S-42~S-43 | `ForbiddenRef` | LookupKey=`subscribe.forbidden` | `ssac/` |
| S-44 | `ForbiddenRef` | LookupKey=`http.forbidden` | `ssac/` |
| S-46 | `NameFormat` | Pattern=`uppercase-start` | `ssac/` |
| S-47 | `NameFormat` | Pattern=`no-dot-prefix` | `ssac/` |
| S-48 | `RefExists` | LookupKey=`SymbolTable.model` | `ssac/` |
| S-49 | `RefExists` | LookupKey=`SymbolTable.method.<Model>` | `ssac/` |
| S-50 | `RefExists` | LookupKey=`OpenAPI.request.<operationId>` | `ssac/` |
| S-54~S-56 | `RefExists` | LookupKey=`OpenAPI.pagination.<operationId>` | `ssac/` |
| S-51 | `CoverageCheck` | LookupKey=`SSaC.requestUsage.<operationId>` | `ssac/` |
| S-53 | `CoverageCheck` | LookupKey=`SSaC.queryUsage` | `ssac/` |
| S-57 | `TypeMatch` | LookupKey=`Func.request.<funcName>` | `ssac/` |
| TM-1~TM-3 | `RefExists` | LookupKey=`OpenAPI.operationId` | `stml/` |
| TM-4 | `RefExists` | LookupKey=`OpenAPI.param.<operationId>` | `stml/` |
| TM-5 | `RefExists` | LookupKey=`OpenAPI.request.<operationId>` | `stml/` |
| TM-6 | `RefExists` | LookupKey=`OpenAPI.response.<operationId>` | `stml/` |
| TM-10 | `RefExists` | LookupKey=`OpenAPI.sort.<operationId>` | `stml/` |
| TM-11 | `RefExists` | LookupKey=`OpenAPI.filter.<operationId>` | `stml/` |

### Defeater 사용

| defeater | 면제 warrant | 조건 | 패키지 |
|----------|-------------|------|--------|
| `IsImplicitVar` | VarDeclared (S-27~S-30) | currentUser, request, query, message | `ssac/` |
| `IsSubscribe` | FieldRequired (HTTP 전용), ForbiddenRef (request→OpenAPI) | @subscribe 함수 | `ssac/` |
| `IsCustomTS` | RefExists (STML bind→OpenAPI) | custom.ts에 함수 존재 | `stml/` |

### 고유 함수 (pkg/rule 미사용)

| 규칙 ID | 함수명 | 패키지 | 설명 |
|---------|--------|--------|------|
| S-11 | `DeleteNoInputs` | `ssac/` | @delete Inputs 없음 WARNING |
| S-25 | `UnknownSeqType` | `ssac/` | 알 수 없는 시퀀스 타입 |
| S-36 | `StaleResponse` | `ssac/` | @put/@delete 후 갱신 없이 @response 사용 |
| S-37 | `FKReferenceGuard` | `ssac/` | FK 참조 @get 후 @empty 가드 필요 |
| S-38~S-41, S-45 | `SubscribeConstraints` | `ssac/` | @subscribe 제약 (파라미터, message struct, @response 금지) |
| S-52 | `QueryUsageMismatch` | `ssac/` | OpenAPI x-pagination ↔ SSaC query 불일치 |
| S-58 | `InvalidErrStatus` | `ssac/` | IANA 미등록 HTTP status |
| C-1 | `ManifestLoad` | `manifest/` | fullend.yaml 로드 실패 |
| D-1 | `SqlcQueryDuplicate` | `ddl/` | sqlc query name 중복 |
| D-2 | `NullableColumn` | `ddl/` | NOT NULL 누락 |
| D-3 | `SentinelMissing` | `ddl/` | FK DEFAULT 0 센티널 누락 |
| O-1 | `PathParamConflict` | `openapi/` | path 파라미터명 충돌 |
| M-1 | `ModelDirEmpty` | `model/` | model/ 디렉토리 비어있음 |
| F-1 | `BuiltinOverride` | `funcspec/` | built-in 패키지명 충돌 |
| H-1 | `DeprecatedFeature` | `hurl/` | .feature 파일 존재 |
| TM-7 | `EachNotArray` | `stml/` | data-each 필드가 배열 아님 |
| TM-8 | `BindNotFound` | `stml/` | data-bind 필드 미발견 (custom.ts도 확인) |
| TM-9 | `PaginateNoExt` | `stml/` | x-pagination 미선언 |
| TM-12 | `ComponentNotFound` | `stml/` | data-component 파일 없음 |
| CT-1~CT-2 | `ContractVerify` | `contract/` | //fullend 디렉티브 검증 |

## 검증 흐름

```
ParseAll() → Fullstack (파싱 결과)
  → SymbolTable 구축 (symbol/)
  → SSOT별 validator 실행:
    - validate/manifest/  → C-1
    - validate/ddl/       → D-1~D-3
    - validate/openapi/   → O-1
    - validate/ssac/      → S-1~S-58 (가장 많음)
    - validate/stml/      → TM-1~TM-12
    - validate/statemachine/ → ST-1
    - validate/rego/      → P-1
    - validate/funcspec/  → F-1
    - validate/hurl/      → H-1
    - validate/model/     → M-1
    - validate/contract/  → CT-1~CT-2
  → 각 validator 내부:
    - rule.Ground 구성 (Lookup, Types, Vars, Flags)
    - Toulmin Graph 구성 (공통 warrant + 고유 warrant + defeater)
    - Graph.Evaluate(claim, ground) per 검증 항목
    - verdict + evidence → ValidationError 변환
```
