✅ 완료

# Phase032: filefunc 100% 준수 — 미분해 파일 전수 분해

## 목표

Phase031에서 분해하지 않은 나머지 21개 파일을 filefunc 규칙으로 분해하여, 4개 대상 패키지의 filefunc 위반을 **F1/F2/A1/A3 = 0**으로 만든다.

설계 원칙:
1. **파일 분리만 수행** — 로직 변경 없음. 함수 시그니처, 호출 관계 동일
2. **같은 패키지 내 분해** — import cycle 위험 없음
3. **기존 기능 무파괴** — `go test ./...` 통과 유지가 전제
4. **Q1(네스팅)/Q3(함수 길이)는 범위 외** — 로직 리팩토링이 필요한 별도 Phase

## 동기

Phase031 이후 남은 미분해 파일:

| 패키지 | 파일 | 줄 | 실제 함수 | 타입 | 설명 |
|---|---|---|---|---|---|
| `internal/gen/gogin/gogin.go` | 605 | 21 | 1 (GoGin) | Gen 디스패처 + 유틸 20개 |
| `internal/gen/gogin/domain.go` | 560 | 10 | 0 | 도메인 분리 코드젠 |
| `internal/gen/gogin/middleware.go` | 324 | 10 | 0 | 인증 패키지 + 미들웨어 코드젠 (내부 템플릿 타입 제외) |
| `internal/gen/gogin/queryopts.go` | 262 | 2 | 0 | QueryOpts 코드젠 (내부 템플릿 타입 제외) |
| `internal/gen/gogin/state.go` | 194 | 5 | 0 | 상태 머신 코드젠 |
| `internal/gen/gogin/main_go.go` | 191 | 5 | 0 | main.go 코드젠 |
| `internal/gen/gogin/server.go` | 167 | 3 | 1 (pathParamInfo) | 서버 구조체 코드젠 |
| `internal/gen/gogin/attach.go` | 165 | 6 | 0 | 핸들러 어태치 코드젠 |
| `internal/gen/gogin/auth.go` | 36 | 1 | 0 | Auth 스텁 코드젠 (이미 1함수, 어노테이션만 추가) |
| `internal/gen/gogin/authz.go` | 31 | 1 | 0 | 인가 코드젠 (이미 1함수, 어노테이션만 추가) |
| `internal/gen/hurl/ddl.go` | 58 | 1 | 1 (ddlFK) | DDL FK 파서 |
| `internal/gen/hurl/hurl_util.go` | 335 | 15 | 0 | Hurl 유틸리티 함수 |
| `internal/orchestrator/gen.go` | 645 | 14 | 0 | Gen 오케스트레이션 |
| `internal/orchestrator/chain.go` | 632 | 17 | 1 (ChainLink) | Chain 오케스트레이션 |
| `internal/orchestrator/detect.go` | 156 | 4 | 3 (SSOTKind+DetectedSSOT+NotDirError)+const | SSOT 감지 |
| `internal/orchestrator/status.go` | 127 | 2 | 1 (StatusLine) | Status 명령 |
| `internal/orchestrator/parsed.go` | 97 | 1 | 0 | ParseAll |
| `internal/orchestrator/sqlc_config.go` | 81 | 2 | 0 | sqlc 설정 생성 |
| `internal/orchestrator/exec.go` | 59 | 2+init | 1 (ExecResult) | CLI 실행 |
| `internal/orchestrator/target_profile.go` | 20 | 1 | 1 (TargetProfile) | 타겟 프로필 |
| `internal/ssac/validator/errors.go` | 26 | 2메서드 | 1 (ValidationError) | ValidationError |

**합계: 21개 파일, ~4,760줄, 실제 함수 125개, 타입 11개**

> **주의**: `middleware.go`와 `queryopts.go`는 백틱 템플릿 문자열 안에 `func`/`type` 선언이 포함되어 있다. 이는 생성될 코드의 일부이며 gogin 패키지의 실제 심볼이 아니다. `grep ^func` 수치와 실제 패키지 심볼 수가 다르므로 분해 시 템플릿 내부 코드를 함수로 오인하지 않도록 주의한다.

## 설계

### 분해 전략

모든 파일에 동일한 filefunc 규칙 적용:
- 1파일 1함수/1타입/1메서드
- `//ff:func feature=xxx type=xxx` + `//ff:what` 어노테이션
- F6 예외: 의미적으로 그룹된 const는 1파일 유지

### 주의 사항

1. **`exec.go`의 `init()` 함수** — `init()`은 PATH 환경변수를 설정하며 `RunExec`이 의존한다. filefunc 규칙상 "init() is only allowed alongside a var or func"이므로, `init()` + `execTimeout` const + `RunExec`은 한 파일에 유지한다. `ExecResult` 타입만 별도 파일로 분리.

2. **`injectFileDirective` 공유 함수** — `attach.go`에 정의, `middleware.go`와 `state.go`에서 호출. 분해 시 별도 파일(`inject_file_directive.go`)로 추출하여 3개 파일이 공유.

3. **어노테이션만 추가하는 파일** — `auth.go`(36줄 1함수), `authz.go`(31줄 1함수), `parsed.go`(97줄 1함수)는 이미 1파일 1개념. `//ff:` 어노테이션만 추가.

4. **`detect.go`의 const 블록** — `SSOTKind` 타입 + 10개 Kind 상수 + `DetectedSSOT` 타입 + `NotDirError` 타입 + 함수 3개. 타입별/함수별 분리 시 const 블록은 `SSOTKind` 타입 파일에 F6 예외로 유지.

5. **`target_profile.go`** — `TargetProfile` 타입 + `DefaultProfile` 함수. F6(unexported type 예외)는 폐지되었으므로 타입과 함수를 분리해야 한다.

6. **`ddl.go`** — `ddlFK` unexported 타입 + `parseDDLFiles` 함수. F6 폐지로 분리 필요.

### 패키지별 분해

#### gen/gogin (10파일 → ~66파일)

| 원본 | 함수 수 | 분해 결과 |
|---|---|---|
| `gogin.go` (605줄) | 21+1타입 | GoGin 타입, internalDir, Generate, collectModels, collectFuncs, transformServiceFiles, transformSource, getExtMap, getStr, getStrSlice, hasBearerScheme, hasDomains, uniqueDomains, collectModelsForDomain, collectFuncsForDomain, domainNeedsAuth, lcFirst, resolveSuccessStatus, httpStatusConst, buildFileToOperationID, collectModelIncludes, collectCursorSpecs 함수별 분리 |
| `domain.go` (560줄) | 10 | transformServiceFilesWithDomains, generateAuthStubWithDomains, generateServerStructWithDomains, domainNeedsJWTSecret, domainNeedsDB, generateDomainHandler, opHasSecurity, convertPathParamsGin, generateCentralServer, generateMainWithDomains 함수별 |
| `middleware.go` (324줄) | 10 | generateAuthPackage, generateIssueToken, generateVerifyToken, generateRefreshToken, generateReexport, generateMiddleware, sortedClaimFields, HashClaimDefs, claimExtractLine, resultAssignLine 함수별 (내부 템플릿 타입은 분해 대상 아님) |
| `queryopts.go` (262줄) | 2 | generateQueryOpts, extractBaseWhere 함수별 (내부 템플릿 타입은 분해 대상 아님) |
| `state.go` (194줄) | 5 | GenerateStateMachines, generateStateMachineSource, generateBoolCanTransition, inferFieldType, inferBoolStates 함수별 |
| `main_go.go` (191줄) | 5 | generateMain, collectSubscribers, hasAuthSequence, buildOwnershipsLiteral, hasPublishSequence 함수별 |
| `server.go` (167줄) | 3+1타입 | generateServerStruct, convertPathParams, ucFirst 함수별 + pathParamInfo 타입 분리 |
| `attach.go` (165줄) | 6 | attachServiceDirectives, attachDirectivesInDir, attachDirectiveToFile, injectFuncDirective, attachTSXDirectives, injectFileDirective(공유→별도 파일) 함수별 |
| `auth.go` (36줄) | 1 | 어노테이션만 추가 (이미 1함수) |
| `authz.go` (31줄) | 1 | 어노테이션만 추가 (이미 1함수) |

#### gen/hurl (2파일 → ~17파일)

| 원본 | 함수 수 | 분해 결과 |
|---|---|---|
| `ddl.go` (58줄) | 1함수+1타입 | parseDDLFiles 함수 + ddlFK 타입 분리 |
| `hurl_util.go` (335줄) | 15 | generateDummyValue, formatDummyValue, generateRequestBody, generateResponseAssertions, resolveSchema, getRequestSchema, getResponseSchema, getSuccessHTTPCode, needsAuth, findTokenJSONPath, inferCaptureField, getExtMap, getStr, getStrSlice, sortStringSlice 함수별 |

#### orchestrator (8파일 → ~49파일)

| 원본 | 함수 수 | 분해 결과 |
|---|---|---|
| `gen.go` (645줄) | 14 | Gen, GenWith, genSqlc, genOpenAPI, genSSaC, injectFuncErrStatusFromParsed, genSTML, genGlue, determineModulePath, genStateMachines, genAuthz, countPolicyRules, genFunc, scanFuncImports 함수별 |
| `chain.go` (632줄) | 17+1타입 | ChainLink 타입, Chain, inferArtifactsDir, traceArtifacts, traceOpenAPI, traceSSaC, traceDDL, tracePolicy, traceStates, traceFuncSpecs, traceHurlScenarios, traceSTML, stmlMatchAttr, findSSaCFile, findDDLTable, grepLine, toSnakeCase, sortedStringKeys 함수별 |
| `detect.go` (156줄) | 4함수+3타입+const | SSOTKind 타입+const(F6 유지), DetectedSSOT 타입, NotDirError 타입, NotDirError.Error 메서드, DetectSSOTs, AllSSOTKinds, KindFromString+kindNames var — 7파일 |
| `status.go` (127줄) | 2함수+1타입 | StatusLine 타입, Status, PrintStatus 함수별 |
| `parsed.go` (97줄) | 1 | 어노테이션만 추가 |
| `sqlc_config.go` (81줄) | 2 | generateSqlcConfig, detectDBEngine 함수별 |
| `exec.go` (59줄) | 2함수+1타입+init | RunExec+init+execTimeout 한 파일 유지, ExecResult 타입만 분리 |
| `target_profile.go` (20줄) | 1함수+1타입 | TargetProfile 타입 + DefaultProfile 함수 분리 |

#### ssac/validator (2파일 → 13파일)

| 원본 | 분해 결과 |
|---|---|
| `errors.go` (26줄) | ValidationError 타입 + Error() 메서드파일 + IsWarning() 메서드파일 |
| `openapi_types.go` (Phase031 잔류) | openAPISpec, openAPIComponents, openAPISchema, openAPIPathItem + operations() 메서드, openAPIOperation, openAPIParameter, openAPIRequestBody, openAPIResponse, openAPIMediaType — 9개 타입/메서드별 분리 (F6 개정으로 타입 그룹 예외 폐지) |

### .ffignore 정리

분해 완료 후 `.ffignore`에서 개별 파일 제외 목록을 삭제하고, `artifacts/`, `pkg/`, `cmd/`, 미적용 패키지 디렉토리만 유지.

## 변경 파일

- `internal/gen/gogin/` — 10개 파일 분해 → ~66개 파일
- `internal/gen/hurl/` — 2개 파일 분해 → ~17개 파일
- `internal/orchestrator/` — 8개 파일 분해 → ~49개 파일
- `internal/ssac/validator/` — errors.go + openapi_types.go 분해 → ~13개 파일
- `.ffignore` — 개별 파일 제외 목록 삭제, 디렉토리 제외만 유지

**예상 파일 수 변화**: 226개 .go → ~349개 .go

## 의존성

- Phase031 완료 후 (✅)
- Phase033(OpenAPI Validation Tags)과 독립

## 검증

1. `go test ./...` — 전체 테스트 통과
2. `go build ./cmd/fullend/` — 빌드 성공
3. `go vet ./...` — 경고 없음
4. `filefunc validate ./internal/gen/gogin/` — F1/F2/A1/A3 = 0
5. `filefunc validate ./internal/gen/hurl/` — F1/F2/A1/A3 = 0
6. `filefunc validate ./internal/orchestrator/` — F1/F2/A1/A3 = 0
7. `filefunc validate ./internal/ssac/validator/` — F1/F2/A1/A3 = 0
8. `fullend validate specs/dummys/zenflow-try05/` — 기존 기능 정상
9. `fullend validate specs/dummys/gigbridge-try02/` — 기존 기능 정상
10. mutest 74건 재실행 — 기존 PASS 유지

## 리스크

- **파일 수 증가** — 226 → ~349개. filefunc 설계 의도와 일치하며, AI 에이전트의 탐색 효율 향상.
- **git diff** — Phase031과 동일. 패키지 단위 커밋으로 관리.
- **.ffignore 간소화** — 개별 파일 제외 제거로 유지보수 부담 감소.
