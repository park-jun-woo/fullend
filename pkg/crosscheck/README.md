# pkg/crosscheck

Toulmin defeats graph 기반 SSOT 간 교차 검증. `pkg/rule` 공통 함수 + 고유 함수로 구성.

## 규칙 → pkg/rule 매핑

### RefExists (~31개)

| 규칙 ID | LookupKey | 설명 |
|---------|-----------|------|
| X-1 | `DDL.column.<table>` | x-sort column → DDL |
| X-3 | `DDL.column.<table>` | x-filter column → DDL |
| X-5 | `DDL.table` | x-include target table → DDL |
| X-9 | `DDL.column.<table>` | OpenAPI property → DDL column (유령) |
| X-15 | `OpenAPI.operationId` | SSaC funcName → OpenAPI |
| X-24 | `States.diagram` | SSaC @state → diagram 존재 |
| X-25 | `States.event.<diagramID>` | @state transition → diagram event |
| X-27 | `DDL.column.<table>` | @state field → DDL column |
| X-23 | `SSaC.funcName` | States transition event → SSaC |
| X-31 | `DDL.table` | @ownership table → DDL |
| X-32 | `DDL.column.<table>` | @ownership column → DDL |
| X-33 | `DDL.table` | @ownership via join table → DDL |
| X-34 | `DDL.column.<table>` | @ownership via join column → DDL |
| X-35 | `OpenAPI.path` | Hurl path → OpenAPI |
| X-36 | `OpenAPI.method.<path>` | Hurl method → OpenAPI |
| X-38 | (함수명 검사) | @call 함수명 소문자 시작 |
| X-39 | `Func.spec` | @call → func 구현 존재 |
| X-43 | `Func.request.<funcName>` | @call Input field → FuncRequest |
| X-49 | `Config.claims` | currentUser.field → claims |
| X-52 | `Config.middleware` | endpoint security → middleware |
| X-53 | `Config.claims.values` | Rego input.claims → claims |
| X-60 | `Authz.checkRequest` | @auth input field → CheckRequest |
| X-63 | `Config.roles` | Rego role → Config roles |
| X-65 | `DDL.check.<table>` | Rego role → DDL CHECK 제약 |
| X-73 | `Config.claims.fields` | JWT @call input → claims |

### CoverageCheck (~9개)

| 규칙 ID | LookupKey | 설명 |
|---------|-----------|------|
| X-16 | `SSaC.funcName` | OpenAPI operationId → SSaC 함수 사용 여부 |
| X-18 | `SSaC.response.<funcName>` | OpenAPI response field → @response 사용 여부 |
| X-20 | `SSaC.response.<funcName>` | OpenAPI field → shorthand @response 사용 여부 |
| X-54 | `Rego.claims` | Config claims → Rego 참조 여부 |
| X-55 | `SSaC.modelRef` | DDL table → SSaC 참조 여부 |
| X-62 | `SSaC.callRef` | func spec → @call 사용 여부 |
| X-64 | `Rego.roles` | Config roles → Rego 사용 여부 |
| X-10 | `OpenAPI.response.<op>` | DDL column → OpenAPI schema 포함 여부 |

### PairMatch (~6개)

| 규칙 ID | LookupKey | 설명 |
|---------|-----------|------|
| X-28 | `Policy.auth` | SSaC @auth (action:resource) → Rego allow |
| X-29 | `SSaC.auth` | Rego allow (action:resource) → SSaC @auth |
| X-50 | `Config.middleware` | OpenAPI securityScheme → Config middleware |
| X-51 | `OpenAPI.security` | Config middleware → OpenAPI securityScheme |
| X-57 | `SSaC.subscribe` | @publish topic → @subscribe |
| X-58 | `SSaC.publish` | @subscribe topic → @publish |

### TypeMatch (~4개)

| 규칙 ID | LookupKey | 설명 |
|---------|-----------|------|
| X-14 | `SQLc.param.<model>` | SSaC input key case ↔ sqlc param |
| X-44 | `Func.request.<funcName>` | @call Input type ↔ Request field type |
| X-69 | `DDL.check.<table>` | DDL CHECK values ↔ OpenAPI enum |

### SchemaMatch (~7개)

| 규칙 ID | LookupKey | 설명 |
|---------|-----------|------|
| X-17 | `OpenAPI.response.<op>` | SSaC @response fields → OpenAPI response |
| X-19 | `OpenAPI.response.<op>` | shorthand @response → OpenAPI response |
| X-42 | `Func.request.<funcName>` | @call Inputs count → FuncRequest fields |
| X-59 | `SSaC.publish.<topic>` | @subscribe message fields → @publish payload |
| X-66 | `OpenAPI.required.<op>` | SSaC used fields → OpenAPI required |
| X-67 | `DDL.varchar.<table>` | DDL VARCHAR(n) → OpenAPI maxLength |
| X-68 | `DDL.check.<table>` | DDL CHECK IN → OpenAPI enum |

### ConfigRequired (~3개)

| 규칙 ID | ConfigKey | 설명 |
|---------|-----------|------|
| X-48 | `backend.auth.claims` | currentUser 사용 → claims 필수 |
| X-52 | `backend.middleware` | endpoint security → middleware 필수 |
| X-56 | `queue.backend` | @publish/@subscribe → queue 설정 필수 |

### Defeater

| defeater | 면제 warrant | Flags 키 |
|----------|-------------|----------|
| `IsSkipped` | 해당 SSOT 전체 | `skipped.<kind>` |
| `IsPkgModel` | RefExists(SSaC→DDL), CoverageCheck(DDL→SSaC) | `pkgModel` |
| `IsDTO` | RefExists(Model→DDL table) | `dto` |
| `IsArchived` | CoverageCheck(DDL→SSaC) | `archived` |
| `IsSensitiveCol` | CoverageCheck(DDL→OpenAPI) | `sensitive` |
| `IsNoSensitive` | RefExists(sensitive pattern) | `nosensitive` |

## 고유 함수 (pkg/rule 미사용)

| 규칙 ID | 함수명 | 설명 |
|---------|--------|------|
| X-2 | `SortColumnNoIndex` | x-sort column에 인덱스 없음 (WARNING) |
| X-4 | `XIncludeInvalidFormat` | x-include 형식 오류 |
| X-6 | `XIncludeNoFK` | x-include FK 제약 없음 (WARNING) |
| X-7 | `CursorMultipleSort` | cursor 모드 x-sort 2개 이상 |
| X-8 | `CursorNonUniqueSort` | cursor sort default UNIQUE 아님 |
| X-11 | `PluralResultType` | @result 타입 복수형 (WARNING) |
| X-12 | `ResultNoDDLTable` | @result 타입 DDL 테이블 없음 (WARNING) |
| X-13 | `InputNotInDDL` | SSaC input DDL 컬럼 없음 (WARNING) |
| X-21 | `ErrStatusNotInOpenAPI` | @empty/@exists/@state/@auth ErrStatus OpenAPI 미정의 |
| X-22 | `ResponseNo2xx` | @response 있는데 OpenAPI 2xx 없음 |
| X-26 | `MissingStateGuard` | 상태 전이 참여하는데 @state 없음 (WARNING) |
| X-30 | `OwnershipNoAnnotation` | resource_owner 참조인데 @ownership 없음 |
| X-37 | `HurlStatusNotDefined` | Hurl status code OpenAPI 미정의 (WARNING) |
| X-40 | `FuncBodyTodo` | func 본체 미구현 (TODO) |
| X-41 | `FuncForbiddenImport` | func I/O 패키지 import 금지 |
| X-45 | `CallResultMissing` | @result 있지만 func Response 없음 |
| X-46 | `CallResultIgnored` | @result 없지만 func Response 있음 (WARNING) |
| X-47 | `CallSourceVarUndefined` | @call arg source 미정의 (WARNING) |
| X-61 | `SensitiveNoAnnotation` | 민감 패턴 컬럼 @sensitive 없음 (WARNING) |
| X-70 | `MaxLengthExceedsVarchar` | OpenAPI maxLength > DDL VARCHAR (WARNING) |
| X-71 | `PasswordNoMinLength` | password 필드 minLength 없음 (WARNING) |
| X-72 | `EmailNoFormat` | email 필드 format 없음 (WARNING) |

## 검증 흐름

```
Fullstack (파싱 결과) + SymbolTable
  → rule.Ground 생성
    - Lookup: DDL columns, OpenAPI operationIds, SSaC funcNames, Config claims/roles 등
    - Pairs: Policy auth pairs, middleware pairs, pub/sub topics
    - Types: DDL column types, Func request field types
    - Schemas: OpenAPI response schemas, Func request fields
    - Config: backend.auth.claims, queue.backend 등
    - Flags: skipped, pkgModel, dto, archived, sensitive 등
  → 규칙 그룹별 Toulmin Graph 구성
    - 공통 warrant: pkg/rule.RefExists, CoverageCheck, PairMatch 등 + backing
    - 고유 warrant: 패키지 내 함수
    - defeater: IsSkipped, IsPkgModel, IsDTO 등 + defeat edges
  → claim별 Graph.Evaluate(claim, ground)
  → verdict + evidence → CrossError 변환
```
