# ZenFlow Report — try06 mod-5 (Add-on #05 Schedule)

## 목표
워크플로우에 cron 스케줄 등록/조회/삭제 기능 추가. `@call` Func을 통해 `pkg/session`을 간접 호출하여 Phase051(package-prefix @model 폐기) 적용을 검증.

## 결과 요약

| 단계 | 결과 |
|---|---|
| fullend validate | **PASS** (ERROR 0, WARNING 12 — 모두 기존) |
| fullend gen | **PASS** (27 service, 10 funcs) |
| go build | **PASS** |
| hurl --test (scenario-schedule) | **FAIL** — SetSchedule 400 |

## 변경 SSOT 목록

| SSOT | 파일 | 변경 |
|---|---|---|
| fullend.yaml | `fullend.yaml` | `session.backend: postgres` 추가 |
| OpenAPI | `api/openapi.yaml` | POST/GET/DELETE `/workflows/{id}/schedule` 3개 엔드포인트 추가 |
| SSaC | `service/schedule/set_schedule.ssac` | @get→@empty→@auth→@call schedule.SetSchedule→@response |
| SSaC | `service/schedule/get_schedule.ssac` | @get→@empty→@auth→@call schedule.GetSchedule→@response |
| SSaC | `service/schedule/delete_schedule.ssac` | @get→@empty→@auth→@call schedule.DeleteSchedule→@response |
| Func | `func/schedule/set_schedule.go` | SetSchedule — 내부적으로 session.Set 호출 |
| Func | `func/schedule/get_schedule.go` | GetSchedule — 내부적으로 session.Get 호출 |
| Func | `func/schedule/delete_schedule.go` | DeleteSchedule — 내부적으로 session.Delete 호출 |
| Policy | `policy/authz.rego` | SetSchedule(admin), GetSchedule(any), DeleteSchedule(admin) |
| Scenario | `tests/scenario-schedule.hurl` | 스케줄 등록→조회→삭제 E2E |

## 발견된 fullend 버그

### BUG030 — POST JSON body를 c.Query()로 읽음

**위치**: `internal/ssac/generator` (codegen)

**증상**: POST `/workflows/{id}/schedule`의 JSON body 필드 `cron`이 `c.Query("cron")`으로 생성됨. 빈 문자열이 전달되어 "invalid cron expression: expected 5 fields, got 0" 에러 발생.

**생성된 코드** (set_schedule.go:23):
```go
cron := c.Query("cron")  // ← 잘못됨
```

**기대 코드**:
```go
var req struct {
    Cron string `json:"cron"`
}
if err := c.ShouldBindJSON(&req); err != nil { ... }
cron := req.Cron
```

**원인 추정**: codegen이 OpenAPI의 `requestBody` + `application/json`이 있는 POST 엔드포인트에서도 일부 필드를 query parameter로 처리. `@call`만 있고 `@post`가 없는 SSaC 함수에서 JSON body 바인딩이 생성되지 않는 것으로 보임.

### BUG031 — session.Init() 미생성

**위치**: `internal/gluegen` (main.go 생성)

**증상**: `fullend.yaml`에 `session.backend: postgres`가 있지만, 생성된 `main.go`에 `session.Init()` 호출이 없음. `pkg/session`의 `defaultModel`이 nil 상태로 남아 runtime panic 발생.

**원인**: Phase051에서 추가한 `session.Init()` 패턴이 gluegen에 아직 반영되지 않음. `queue.Init()`처럼 main.go에 자동 주입해야 함.

## 설계 결정

### session Key 타입 불일치 우회

`session.Set/Get/Delete`의 `Key`가 `string`인데, SSaC에서 `wf.ID`(int64)를 직접 전달하면 타입 불일치. SSaC에는 타입 변환 기능이 없으므로, `schedule` 패키지 함수에서 `fmt.Sprintf("schedule:%d", WorkflowID)`로 키를 구성하고 내부적으로 `pkg/session`을 호출하는 방식으로 우회.

## fullend chain 결과

```
── Feature Chain: SetSchedule ──
  OpenAPI    api/openapi.yaml:987    POST /workflows/{id}/schedule
  SSaC       service/schedule/set_schedule.ssac:10
  DDL        db/workflows.sql:1      CREATE TABLE workflows
  Rego       policy/authz.rego:5     resource: workflow
  FuncSpec   func/schedule/set_schedule.go:14

── Feature Chain: GetSchedule ──
  OpenAPI    api/openapi.yaml:1039   GET /workflows/{id}/schedule
  SSaC       service/schedule/get_schedule.ssac:10
  DDL        db/workflows.sql:1      CREATE TABLE workflows
  Rego       policy/authz.rego:5     resource: workflow
  FuncSpec   func/schedule/get_schedule.go:12

── Feature Chain: DeleteSchedule ──
  OpenAPI    api/openapi.yaml:1074   DELETE /workflows/{id}/schedule
  SSaC       service/schedule/delete_schedule.ssac:10
  DDL        db/workflows.sql:1      CREATE TABLE workflows
  Rego       policy/authz.rego:5     resource: workflow
  FuncSpec   func/schedule/delete_schedule.go:12
```

## 결론

Phase051(package-prefix @model 폐기)의 SSOT 수준 적용은 성공. `@call` Func을 통해 `pkg/session`을 간접 호출하는 패턴으로 validate + gen + build까지 통과. Hurl E2E는 codegen 버그 2건(BUG030, BUG031)으로 실패. 코드 수정 없이 fullend 도구 수정이 필요함.
